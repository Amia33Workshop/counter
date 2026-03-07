{ pkgs, ... }: {
  channel = "unstable";
  packages = [
    pkgs.go
    pkgs.gopls
    pkgs.gcc
    pkgs.gnupg
    pkgs.openssh
  ];
  idx = {
    extensions = [
      "mhutchie.git-graph"
      "oderwat.indent-rainbow"
      "esbenp.prettier-vscode"
      "google.gemini-cli-vscode-ide-companion"
      "golang.Go"
    ];
    workspace = {
      onCreate = {
        download = "go mod download";
        tidy = "go mod tidy";
        verify = "go mod verify";
        default.openFiles = [ "templates/index.html" ];
      };
      onStart = { default.openFiles = [ "templates/index.html" ]; };
    };
  };
}
