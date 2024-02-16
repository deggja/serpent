class serpent < Formula
    desc "Play snake and wreck havoc inside your Kubernetes cluster at the same time"
    homepage "https://github.com/deggja/serpent"
    url "URL to the tarball of the release (e.g., https://github.com/user/repo/archive/v1.0.0.tar.gz)"
    sha256 "SHA-256 of the tarball"
  
    depends_on "dependency_name" => :optional
  
    def install
      system "make", "install", "PREFIX=#{prefix}"
    end
  
    test do
      system "#{bin}/your_app_name", "--version"
    end
  end
  