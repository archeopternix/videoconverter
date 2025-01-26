# Video Converter

A simple video converter application built with Go and Fyne.

## Features

- Drag and drop video files for conversion.
- Select video files via a file dialog.
- Configure video conversion settings.
- Save converted videos to a selected directory.
- Progress bar to show conversion progress.

## Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/archeopternix/videoconverter.git
    ```
2. Navigate to the project directory:
    ```sh
    cd videoconverter
    ```
3. Install dependencies:
    ```sh
    go get ./...
    ```
4. Build the application:
    ```sh
    go build
    ```

## Usage

1. Run the application:
    ```sh
    ./videoconverter
    ```
2. Drag and drop video files into the application window or use the file dialog to select files.
3. Configure the video conversion settings using the "Settings" button.
4. Select the destination folder for converted videos.
5. Click the "Convert" button to start the conversion process.

## Dependencies

- [Fyne](https://fyne.io/) - A cross-platform GUI toolkit for Go.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

Feel free to modify or extend this README to better suit the needs of your project.