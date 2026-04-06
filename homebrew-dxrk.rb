# Documentation: https://docs.brew.sh/Formula-Cookbook
class Dxrk < Formula
  desc 'Dxrk AI — AI coding agent ecosystem configurator'
  homepage 'https://github.com/Dxrk777/Dxrk-Hex'
  version '1.3.0'
  license 'MIT'

  on_macos do
    on_intel do
      url 'https://github.com/Dxrk777/Dxrk-Hex/releases/download/v1.3.0/dxrk_1.3.0_darwin_amd64.tar.gz'
      sha256 '5d9a38f4fe771b993ef3327aa9c6a24ce6a35558760c6c6b4168a8764cd6e7ad'
    end
    on_arm do
      url 'https://github.com/Dxrk777/Dxrk-Hex/releases/download/v1.3.0/dxrk_1.3.0_darwin_arm64.tar.gz'
      sha256 '2de2dcdcce2d2a524da82c0011e01f0e2b6401bd66246b3356172c3fed50efdf'
    end
  end

  on_linux do
    on_intel do
      url 'https://github.com/Dxrk777/Dxrk-Hex/releases/download/v1.3.0/dxrk_1.3.0_linux_amd64.tar.gz'
      sha256 '461287e912442a823338334685487d24e096a14ec4c217cfde922c1522a3fa12'
    end
    on_arm do
      url 'https://github.com/Dxrk777/Dxrk-Hex/releases/download/v1.3.0/dxrk_1.3.0_linux_arm64.tar.gz'
      sha256 'f6ef9b76648f2d4567f559d77828633f42a32bec51dfb813695c0d01df9ce21e'
    end
  end

  def install
    bin.install 'dxrk'
  end

  test do
    system "#{bin}/dxrk", 'version'
  end
end
