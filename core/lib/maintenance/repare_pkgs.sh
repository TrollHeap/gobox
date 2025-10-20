#!/usr/bin/env bash
set -euo pipefail

source "$DIR_ROOT/lib/utils/echo_status.sh"

repare_pkgs(){
    echo_status "Mise à jour des index de paquets"
    sudo apt update && echo_status_ok "Update réussi"|| echo_status_error "Échec apt update"

    echo_status "Réparation des paquets cassés"
    sudo apt --fix-broken install -y && echo_status_ok "Réparation des paquets cassés."|| echo_status_error "Échec fix-broken"

    echo_status "Nettoyage des paquets obsolètes"
    sudo apt autoclean -y && echo_status_ok "Nettoyage des paquets."|| echo_status_error "Échec autoclean"

    echo_status "Suppression des dépendances inutiles"
    sudo apt autoremove -y && echo_status_ok "Suppression des dépendances"|| echo_status_error "Échec autoremove"
}

