{
  pkgs ? import <nixpkgs> { },
}:

pkgs.mkShell {
  buildInputs = with pkgs; [
    gcc
    cmake
    ninja
    boost
    pkg-config
    gdb
  ];

  shellHook = ''
    echo "C++ dev environment loaded"
    jetbrains-toolbox
  '';
}
