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
      onCreate = {
        download = "go mod download";
        tidy = "go mod tidy";
        verify = "go mod verify";
        default.openFiles = [ "templates/index.html" ];
      };
      onStart = { };
    };
  };
}
