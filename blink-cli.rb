# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class BlinkCli < Formula
  desc "Awesome CLI for the blink ops platform"
  homepage "https://www.blinkops.com/"
  version "0.3.12"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/blinkops/blink-go-cli/releases/download/v0.3.12/blink-cli_0.3.12_darwin_arm64.tar.gz"
      sha256 "98a6bdc481ad83c96a93a56bcdece22c4c5171c940b5f480d9d3cd580e8284e1"

      def install
        bin.install "blink-cli"
      end
    end
    if Hardware::CPU.intel?
      url "https://github.com/blinkops/blink-go-cli/releases/download/v0.3.12/blink-cli_0.3.12_darwin_amd64.tar.gz"
      sha256 "61121112659d63a13ded5adfd6f82b2ab7ce9a4fb38d26cd567afbab0a5ac74b"

      def install
        bin.install "blink-cli"
      end
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "https://github.com/blinkops/blink-go-cli/releases/download/v0.3.12/blink-cli_0.3.12_linux_amd64.tar.gz"
      sha256 "3ce9450a88374d80d9f16b2c571e0d336286b2651b25950a73cbd27c2e291c74"

      def install
        bin.install "blink-cli"
      end
    end
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/blinkops/blink-go-cli/releases/download/v0.3.12/blink-cli_0.3.12_linux_arm64.tar.gz"
      sha256 "d427150734977fe1dcd8207554bf55e8270995f7d13f37b9c93d3cdb24e51a4c"

      def install
        bin.install "blink-cli"
      end
    end
  end
end
