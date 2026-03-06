{ pkgs, ... }: {
  channel = "unstable";
  packages = [
    pkgs.go
    pkgs.gopls
    pkgs.gcc
    pkgs.gnupg
    pkgs.openssh
  ];
  env = { };
  idx = {
    extensions = [
      "mhutchie.git-graph"
      "oderwat.indent-rainbow"
      "esbenp.prettier-vscode"
      "google.gemini-cli-vscode-ide-companion"
      "golang.Go"
    ];
    previews = {
      enable = false;
    };
    workspace = {
      onCreate = { };
      onStart = { };
    };
  };
}
