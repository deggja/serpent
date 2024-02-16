class Netfetch < Formula
  desc "Play snake in your terminal and wreck havoc to your Kubernetes cluster. Lol."
  homepage "https://github.com/deggja/serpent"

  if OS.mac?
    url ""
    sha256 ""
  elsif OS.linux?
    url ""
    sha256 ""
  end

  def install
    bin.install "serpent"
  end

  test do
    system "#{bin}/serpent", "version"
  end
end
