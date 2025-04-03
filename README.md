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

### Prerequisites

- Go 1.18 or higher

### Building from Source

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
   - `KeyPrefix`: Required prefix for encryption keys (default: lh-)
   - `KeySuffix`: Required suffix for encryption keys (default: u)
   - `RequiredChars`: Characters that must be present in the key (default: !@#$%^&*)
   - `MinSpecialChars`: Minimum number of special characters required (default: 2)
   - `KeyContain`: String that must be contained in the key (default: key)

This way, only you know the exact rules for valid encryption keys, making it much harder for others to guess your keys even if they have access to your encrypted data.

## Usage

### Storing a New API Key

```bash
./lhkeymanager
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
Enter environment variable name (with suffix): OPENAI_API_KEY_PROD

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

## Acknowledgements

- This project uses AES-256-GCM encryption for secure key storage
- Inspired by the need for secure API key management in development environments

