# Documentation: https://docs.brew.sh/Formula-Cookbook
class Dxrk < Formula
  desc "Dxrk Hex — AI coding agent ecosystem configurator"
  homepage "https://github.com/Dxrk777/Dxrk-Hex"
  version "1.17.0"
  license "MIT"

  on_macos do
    on_intel do
      url "https://github.com/Dxrk777/Dxrk-Hex/releases/download/v1.17.0/dxrk_1.17.0_darwin_amd64.tar.gz"
      sha256 "REPLACE_WITH_ACTUAL_SHA256"
    end
    on_arm do
      url "https://github.com/Dxrk777/Dxrk-Hex/releases/download/v1.17.0/dxrk_1.17.0_darwin_arm64.tar.gz"
      sha256 "REPLACE_WITH_ACTUAL_SHA256"
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/Dxrk777/Dxrk-Hex/releases/download/v1.17.0/dxrk_1.17.0_linux_amd64.tar.gz"
      sha256 "REPLACE_WITH_ACTUAL_SHA256"
    end
    on_arm do
      url "https://github.com/Dxrk777/Dxrk-Hex/releases/download/v1.17.0/dxrk_1.17.0_linux_arm64.tar.gz"
      sha256 "REPLACE_WITH_ACTUAL_SHA256"
    end
  end

  def install
    bin.install "dxrk"
  end

  test do
    system "#{bin}/dxrk", "version"
  end
end
