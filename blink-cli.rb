# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class BlinkCli < Formula
  desc "Awesome CLI for the blink ops platform"
  homepage "https://www.blinkops.com/"
  version "0.3.10"
  license "MIT"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/blinkops/blink-go-cli/releases/download/v0.3.10/blink-cli_0.3.10_darwin_amd64.tar.gz"
      sha256 "edc4103da96f7bf519c395065354e78dc2adb4a3315032c25bc19668e9e22e34"

      def install
        bin.install "blink-cli"
      end
    end
    if Hardware::CPU.arm?
      url "https://github.com/blinkops/blink-go-cli/releases/download/v0.3.10/blink-cli_0.3.10_darwin_arm64.tar.gz"
      sha256 "0cc159cffe4dc46b4164676ba0b1a42f3ec7ad4b1c31e68f77b9877b151493a6"

      def install
        bin.install "blink-cli"
      end
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "https://github.com/blinkops/blink-go-cli/releases/download/v0.3.10/blink-cli_0.3.10_linux_amd64.tar.gz"
      sha256 "368fd1b0b12a826e68dcc7026cf22704298ea72b6e3227a04ed3429dde555930"

      def install
        bin.install "blink-cli"
      end
    end
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/blinkops/blink-go-cli/releases/download/v0.3.10/blink-cli_0.3.10_linux_arm64.tar.gz"
      sha256 "40c0eceacd74765c5851557e38b9de5704e61e7d382991aa24e10aa0d32f4c89"

      def install
        bin.install "blink-cli"
      end
    end
  end
end
