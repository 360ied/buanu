{ pkgs ? import <nixos> { } }:
pkgs.mkShell { nativeBuildInputs = [ pkgs.go_1_18 pkgs.gotools ]; }
