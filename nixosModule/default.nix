{ config, lib, pkgs, self, ... }:
let
  cfg = config.services.TorrentHelperBot;
in {
  options.services.TorrentHelperBot = {
    enable = lib.mkEnableOption "Enable torrent helper bot service";
    package = lib.mkOption {
      type = lib.types.package;
      default = self.packages.${pkgs.stdenv.hostPlatform.system}.default;
      description = "Package for torrent helper bot";
    };
    allowedUsers = lib.mkOption {
      type = lib.types.listOf lib.types.str;
      default = [];
      description = "lib of allowed users to talk";
    };
    telegramBotTokenPath = lib.mkOption {
      type = lib.types.path;
      description = "path to a file that contains the secret telegram bot token";
    };
    transmissionRpcUrl = lib.mkOption {
      type = lib.types.str;
      default = "http://localhost:9091/transmission/rpc";
      description = "Url of the transmission rpc.";
    };
  };

  config = lib.mkIf cfg.enable {
    systemd.services."torrent-helper-bot" = {
      environment = {
        TORRENT_HELPER_BOT_ALLOWED_USERS = lib.escapeShellArgs cfg.allowedUsers;
        TORRENT_HELPER_BOT_TELEGRAM_BOT_TOKEN = cfg.telegramBotTokenPath;
        TORRENT_HELPER_BOT_TRANSMISSION_RPC_URL = cfg.transmissionRpcUrl;
      };
      wantedBy = [ "multi-user.target" ];
      after = [ "network-online.target" ];
      wants = [ "network-oneline.target" ];
      serviceConfig = {
        ExecStart = "${cfg.package}/bin/TorrentHelperBot";
        Restart = "on-failure";
        RestartSec = 5;
      };

      DynamicUser = true;
      StateDirectory = "TorrentHelperBot";
      RuntimeDirectory = "TorrentHelperBot";
      WorkingDirecctory = "/var/lib/TorrentHelperBot";

      NoNewPrivileges = true;
      ProtectSystem = "strict";
      ProtectHome = true;
      PrivateTmp = true;
      PrivateDevices = true;
      ProtectKernelTunables = true;
      ProtectKernelModules = true;
      ProtectControlGroups = true;
      RestrictAddressFamilies = [ "AF_INET" ];
      RestrictNamespaces = true;
      LockPersonality = true;
      MemoryDenyWriteExecute = true;
      SystemCallArchitectures = "native";
      SystemCallFilter = [ "@system-service" "~@privileged" "~@resources" ];
    };
  };
}
