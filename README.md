# Akeyless JSON Secret Importer

A command-line application written in Go that reads JSON secrets from a directory and its subdirectories, then makes API calls to the Akeyless API to create and update secrets.

## Features

- Recursively scans a directory for JSON files containing secrets.
- Decodes base64-encoded secrets.
- Makes API calls to the Akeyless API to create and update secrets.

## Environment Variables

- `AKEYLESS_TOKEN`: Required. The token used for Akeyless API calls. If this environment variable is not set, the application will prompt for it at runtime.
- `AKEYLESS_IMPORT_STARTING_PATH`: Optional. The directory from which the application should start scanning for JSON files. If not set, the application will start from the current working directory (`"."`).
- `AKEYLESS_SECRET_NAME_PREFIX`: Optional. A prefix to prepend to the secret name for every API call.
- `AKEYLESS_API_GW_URL`: Optional. The URL to the Akeyless API Gateway. If not set, the application will default to `"https://api.akeyless.io"`.

## Defaults

- `AKEYLESS_IMPORT_STARTING_PATH`: `"."`
- `AKEYLESS_API_GW_URL`: `"https://api.akeyless.io"`

## Usage

1. Set the necessary environment variables.
2. Run the application.

```sh
./json-secret-importer
```

## Build

The application can be built for Linux, macOS (both amd64 and arm64), and Windows using GoReleaser. Check the `goreleaser.yml` configuration file for more details.

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[MIT](https://choosealicense.com/licenses/mit/)
