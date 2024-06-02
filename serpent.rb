class Serpent < Formula
  desc "Play snake in your terminal and wreck havoc to your Kubernetes cluster. Lol."
  homepage "https://github.com/deggja/serpent"

  if OS.mac?
    url "https://github.com/deggja/serpent/releases/download/3.1.0/serpent_3.1.0_darwin_amd64.tar.gz"
    sha256 "7aa204a21c8fec041f8f18e9b686cf9bc115e3738e60a764baf6cfab92a663e2"
  elsif OS.linux?
    url "https://github.com/deggja/serpent/releases/download/3.1.0/serpent_3.1.0_linux_amd64.tar.gz"
    sha256 "615de4c69ee013301937d9a47132b4d3221395d71b4f93d39550b23219e3eb1b"
  end

  def install
    bin.install "serpent"
  end

  test do
    system "#{bin}/serpent", "version"
  end
end
