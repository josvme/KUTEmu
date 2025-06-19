{
  description = "RiscV Cross Compiler";

  inputs = {
    flake-parts.url = "github:hercules-ci/flake-parts";
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      imports = [
        # To import a flake module
        # 1. Add foo to inputs
        # 2. Add foo as a parameter to the outputs function
        # 3. Add here: foo.flakeModule

      ];
      systems = [ "x86_64-linux" ];
      perSystem = { config, self', inputs', pkgs, system, ... }: {
        packages = {
          # callPackage resolve dependencies for the package.nix inputs, like stdenv based on the used pkgs
          default = pkgs.callPackage ./package.nix { };
          clang = pkgs.callPackage ./package.nix {stdenv = pkgs.clangStdenv;};
          riscv-bare = pkgs.pkgsCross.riscv32-embedded.callPackage ./package-riscv-bare.nix { };
        };
        # devShells.default describes the default shell with C++, cmake
        devShells = {
          default = config.packages.default;
          clang = config.packages.clang;
          riscv = config.packages.riscv-bare;
        };
      };
      flake = {
        # The usual flake attributes can be defined here, including system-
        # agnostic ones like nixosModule and system-enumerating ones, although
        # those are more easily expressed in perSystem.

      };
    };
}
