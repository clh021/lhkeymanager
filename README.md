# LH Key Manager

[README in Chinese](README_zh.md)

A secure tool for managing API keys, encrypting them in `.env` files, and loading them into new bash sessions when needed.

## Features

- Securely encrypt and store API keys in `.env` files
- Load encrypted keys into new bash sessions
- Environment variables only exist in the new bash session and are automatically cleared when the session ends
- Uses AES-256-GCM encryption algorithm for key protection
- Interactive input to prevent sensitive information from being recorded in bash history

## Installation

### Download from GitHub Releases

The easiest way to install LH Key Manager is to download a pre-built binary from the [GitHub Releases](https://github.com/clh021/lhkeymanager/releases) page.

1. Go to the [Releases](https://github.com/clh021/lhkeymanager/releases) page
2. Download the appropriate version for your operating system:
   - For Linux: `lhkeymanager-vX.Y.Z-linux-amd64.tar.gz` or `lhkeymanager-vX.Y.Z-linux-arm64.tar.gz`
   - For macOS: `lhkeymanager-vX.Y.Z-darwin-amd64.tar.gz` or `lhkeymanager-vX.Y.Z-darwin-arm64.tar.gz`
   - For Windows: `lhkeymanager-vX.Y.Z-windows-amd64.zip`
3. Extract the archive and make the binary executable (Linux/macOS):
   ```bash
   # For Linux/macOS
   tar -xzf lhkeymanager-vX.Y.Z-linux-amd64.tar.gz
   chmod +x lhkeymanager-vX.Y.Z-linux-amd64
   # Optionally move to a directory in your PATH
   sudo mv lhkeymanager-vX.Y.Z-linux-amd64 /usr/local/bin/lhkeymanager
   ```

### Building from Source

If you prefer to build from source:

#### Prerequisites

- Go 1.18 or higher

```bash
# Clone the repository
git clone https://github.com/clh021/lhkeymanager.git
cd lhkeymanager

# Build using the build script (recommended)
./build.sh  # English version
# or
./build_zh.sh  # Chinese version

# Or build manually
go build -o lhkeymanager
```

### Customizing Security Rules

For enhanced security, you can customize the encryption key validation rules during the build process:

1. Run the build script and choose to customize security rules when prompted:

   ```bash
   ./build.sh
   ```

2. The script will ask you to configure the following security rules:
   - `MinKeyLength`: Minimum length for encryption keys (default: 16)
   - `KeyPrefix`: Required prefix for encryption keys (default: lh-, enter 'empty' for no prefix)
   - `KeySuffix`: Required suffix for encryption keys (default: u, enter 'empty' for no suffix)
   - `RequiredChars`: Characters that must be present in the key (default: !@#$%^&\*, enter 'empty' for no special character requirements)
   - `MinSpecialChars`: Minimum number of special characters required (default: 2)
   - `KeyContain`: String that must be contained in the key (default: key, enter 'empty' for no content requirements)

This way, only you know the exact rules for valid encryption keys, making it much harder for others to guess your keys even if they have access to your encrypted data.

## Usage

### Storing a New API Key

```bash
./lhkeymanager store [file_path]
```

Select option `1`, then follow the prompts to enter your encryption key and API key.

### Loading Keys into a New Bash Session

```bash
./lhkeymanager
```

Select option `2`, enter your encryption key, and the tool will start a new bash session with the environment variables set.

## Security Considerations

- The `.env` file permissions are automatically set to 600 (readable and writable only by the owner)
- The encryption key is never stored and must be manually entered each time
- Environment variables only exist in the new bash session and are cleared when the session ends
- Temporary files are securely deleted after use

## Examples

### Storing an OpenAI API Key

```
$ ./lhkeymanager
Please select an operation:
1. Store a new API key in the .env file
2. Load keys from the .env file into a new bash session
Enter your choice (1/2): 1
Enter encryption key: [input not shown]
Enter API key to encrypt: [input not shown]
Enter environment variable name: OPENAI_API_KEY

Encryption result: enc:AES256:AAECAwQFBgcICQoLDA0ODxAREhMUFRYXGBkaGxwdHh8gISIjJCUmJygpKissLS4vMDEyMzQ1Njc4OTo7PD0+P0BBQkNERUZHSElKS0xNTk9QUVJTVFVWV1hZWltcXV5fYGFiY2RlZmdo
Successfully saved to .env file
```

### Loading Keys into a New Bash Session

```
$ ./lhkeymanager
Please select an operation:
1. Store a new API key in the .env file
2. Load keys from the .env file into a new bash session
Enter your choice (1/2): 2
Enter encryption key: [input not shown]
Environment variable set: OPENAI_API_KEY

Starting a new bash session with environment variables...
$ echo $OPENAI_API_KEY
sk-your-api-key
$ exit
Bash session ended, environment variables cleared
```

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Creating a New Release

To create a new release:

1. Use the provided script:

   ```bash
   ./create-release.sh v1.0.0
   ```

   This will create and push a new tag, which will trigger the GitHub Actions workflow to build and publish the release.

2. Alternatively, you can manually create and push a tag:

   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

3. You can also trigger a manual release from the GitHub Actions tab by selecting the "Release" workflow and clicking "Run workflow". This allows you to specify custom security rules for the build.

## Acknowledgements

- This project uses AES-256-GCM encryption for secure key storage
- Inspired by the need for secure API key management in development environments
