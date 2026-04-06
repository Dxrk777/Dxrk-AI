# Documentation: https://docs.brew.sh/Formula-Cookbook
class Dxrk < Formula
  desc 'Dxrk AI — AI coding agent ecosystem configurator'
  homepage 'https://github.com/Dxrk777/Dxrk-Hex'
  version '1.1.0'
  license 'MIT'

  on_macos do
    on_intel do
      url 'https://github.com/Dxrk777/Dxrk-Hex/releases/download/v1.3.0/dxrk_1.1.0_darwin_amd64.tar.gz'
      sha256 '693f15aa41549406806b7d8fece345bf13e12b2761a4fce01e037912f082982f'
    end
    on_arm do
      url 'https://github.com/Dxrk777/Dxrk-Hex/releases/download/v1.3.0/dxrk_1.1.0_darwin_arm64.tar.gz'
      sha256 'dd83b345021a7a3d89e82ebb8ecd9f5131e1954051cbcd3d2e8891f07c234bdd'
    end
  end

  on_linux do
    on_intel do
      url 'https://github.com/Dxrk777/Dxrk-Hex/releases/download/v1.3.0/dxrk_1.1.0_linux_amd64.tar.gz'
      sha256 '6b2ee33f2f358c2526cce4104e4a6047e9972b9621b4fb4d6d3849e9c588e6a2'
    end
    on_arm do
      url 'https://github.com/Dxrk777/Dxrk-Hex/releases/download/v1.3.0/dxrk_1.1.0_linux_arm64.tar.gz'
      sha256 '73abf2829a8c2c4f7360cf641f0f69cea04f576bb8929ff130a72aed3ffd321c'
    end
  end

  def install
    bin.install 'dxrk'
  end

  test do
    system "#{bin}/dxrk", 'version'
  end
end
