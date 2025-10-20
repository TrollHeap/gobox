#!/usr/bin/env bash
set -euo pipefail

CORE_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
export CORE_DIR
source "$CORE_DIR/etc/config/find_project_root.sh" || echo "Error sourcing find_project_root.sh"
source "$CORE_DIR/etc/config/path.env"           || echo "Error sourcing path.env"
source "$CORE_DIR/lib/ui/echo_status.sh"             || echo "Error sourcing echo_status.sh"

declare -A ACTION_DESC=(
    [PRINT]="Outils d'impression/rapport"
    [MAJ]="Install/Mise à jour système"
    [TEST_HW]="Test de matériel v2"
    [CLONE]="Clonage de partitions/disques"
    [SHRED]="Shred disque dur"
    [CLEAN]="Nettoyage des fichiers inutiles"
)

list_actions() {
    local folder="bin"
    while IFS= read -r -d '' file; do
        name="${file##*/}"      # basename
        name="${name%.sh}"      # retire la dernière occurrence de .sh
        printf '%s\n' "$name"
    done < <(find "$folder" -type f -name "*.sh" -print0)
}

build_menu_items() {
    local menu=()
    for action in $(list_actions); do
        local desc="${ACTION_DESC[$action]:-}"
        menu+=("$action" "$desc")
    done
    menu+=("QUITTER" "Sortir de l'outil")
    printf '%s\n' "${menu[@]}"
}

menu_tui() {
    whiptail --title "Bienvenue sur la Toolbox" --msgbox "Bienvenue sur la Toolbox" 8 40

    while :; do
        # Rafraîchit dynamiquement le menu
        IFS=$'\n' read -r -d '' -a ITEMS < <(build_menu_items && printf '\0')

        ACTION=$(whiptail \
                --title "Sélectionnez une action" \
                --menu "Choisissez une action à exécuter :" 20 70 12 \
                "${ITEMS[@]}" \
                3>&1 1>&2 2>&3
        ) || { echo "Action annulée." >&2; exit 1; }

        if ! [[ "$ACTION" =~ ^[A-Za-z0-9_-]+$ ]]; then
            echo_status_error "Action invalide : $ACTION"
        fi

        [[ "$ACTION" == QUITTER ]] && exit 0

        SCRIPT="$BIN_DIR/${ACTION}.sh"
        if [[ -f "$SCRIPT" && -x "$SCRIPT" ]]; then
            bash "$SCRIPT"
        elif [[ -f "$SCRIPT" ]]; then
            bash "$SCRIPT"
        else
            echo_status_error "Action inconnue ou script absent : $ACTION"
            sleep 2
        fi

        whiptail --title "Action terminée" --msgbox "Action $ACTION terminée." 6 40
    done
}

menu_tui
