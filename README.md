<a href="https://www.buymeacoffee.com/ufukvarol4"><img src="https://img.buymeacoffee.com/button-api/?text=Buy me a coffee&emoji=&slug=ufukvarol4&button_colour=40DCA5&font_colour=ffffff&font_family=Cookie&outline_colour=000000&coffee_colour=FFDD00" /></a>

# ISO8583 Payment Simulator

ISO8583 Payment Simulator is an open-source project designed to handle ISO8583 payment requests from POS terminals and provide dummy responses. This tool is perfect for developers and testers working on payment systems, enabling them to simulate and test transactions without needing a live payment gateway.

## Features

- **Real-Time ISO8583 Handling:** Seamlessly processes ISO8583 payment requests from POS terminals.
- **Customizable Responses:** Easily configure and customize dummy responses for different test scenarios.
- **Detailed Logging:** Comprehensive logging to monitor and analyze the flow of requests and responses.
- **Scalability:** Designed to handle a high volume of transactions, making it suitable for large-scale testing environments.
- **User-Friendly Interface:** Simple and intuitive interface for easy configuration and management.

## Getting Started

### Prerequisites

- Go 1.16+ installed on your machine.

### Installation

1. **Clone the Repository:**
   ```bash
   git clone https://github.com/thearistotlemethod/MockIso8583Host.git
   ```
2. **Navigate to the project directory:**
   ```bash
   cd MockIso8583Host
   ```
3. **Build the project:**
   ```bash
   go build
   ```

### Usage

1. **Run the Simulator:**
   ```bash
   ./MockIso8583Host
   ```
2. **Configure your POS terminal to point to the simulator's IP address and port.**

## Contributing

We welcome contributions from the community! If you'd like to contribute, please fork the repository and submit a pull request.

## License

This project is licensed under the Apache-2.0 License - see the [LICENSE](LICENSE) file for details.

## Acknowledgements

Thanks to all the contributors who have made this project possible.
