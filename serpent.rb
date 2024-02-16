class Netfetch < Formula
  desc "Play snake in your terminal and wreck havoc to your Kubernetes cluster. Lol."
  homepage "https://github.com/deggja/serpent"

  if OS.mac?
    url "https://github.com/deggja/serpent/releases/download/2.0.0/serpent_2.0.0_darwin_amd64.tar.gz"
    sha256 "88303dd3b51ebc936169d591ead716f74627388e9944328e86bd2a72a4c9870d"
  elsif OS.linux?
    url "https://github.com/deggja/serpent/releases/download/2.0.0/serpent_2.0.0_linux_amd64.tar.gz"
    sha256 "7c508646a76a9337b7b03370e93a49caa829b57fffa362fee08ec0af6f442225"
  end

  def install
    bin.install "serpent"
  end

  test do
    system "#{bin}/serpent", "version"
  end
end
