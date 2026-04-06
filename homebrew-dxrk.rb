# Documentation: https://docs.brew.sh/Formula-Cookbook
class Dxrk < Formula
  desc 'Dxrk AI — AI coding agent ecosystem configurator'
  homepage 'https://github.com/Dxrk777/Dxrk'
  version '1.3.0'
  license 'MIT'

  on_macos do
    on_intel do
      url 'https://github.com/Dxrk777/Dxrk/releases/download/v1.3.0/dxrk_1.3.0_darwin_amd64.tar.gz'
      sha256 'ebc743634707c6c948c108ddfbed68a17eea3f30b15978377518fbbf85b4ea00'
    end
    on_arm do
      url 'https://github.com/Dxrk777/Dxrk/releases/download/v1.3.0/dxrk_1.3.0_darwin_arm64.tar.gz'
      sha256 '31d53fe9b793da8a8752d78b7d31a5c51188db2fe7939c3ce86f97336a46f463'
    end
  end

  on_linux do
    on_intel do
      url 'https://github.com/Dxrk777/Dxrk/releases/download/v1.3.0/dxrk_1.3.0_linux_amd64.tar.gz'
      sha256 '789fd9ab4e31ef3e303c1e109924257b75012d48c321e33eed73982a77bbdc19'
    end
    on_arm do
      url 'https://github.com/Dxrk777/Dxrk/releases/download/v1.3.0/dxrk_1.3.0_linux_arm64.tar.gz'
      sha256 '1365a1e5d0ffce9158e20b4dd9e6b7eb014d00cc9521ac4e7847a09cc8775d91'
    end
  end

  def install
    bin.install 'dxrk'
  end

  test do
    system "#{bin}/dxrk", 'version'
  end
end
