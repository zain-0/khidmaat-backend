# Khidmaat Backend

## Overview
The Khidmaat Backend is a Go-based server application designed to manage and facilitate various services. It provides APIs for handling users, devices, hospitals, and medical records. The project is structured to ensure scalability, maintainability, and ease of development.

## Project Structure
The project follows a modular structure:

- **`main.go`**: The entry point of the application.
- **`config/`**: Contains configuration files, such as database setup (`db.go`).
- **`controllers/`**: Houses the logic for handling requests and responses for different entities like users, devices, hospitals, and ECG data.
- **`models/`**: Defines the data models for entities such as users, devices, hospitals, and medical records.
- **`routers/`**: Manages the routing of API endpoints for different modules.
- **`utils/`**: Contains utility functions and helpers, such as signal processing helpers (`denoise-helper.go`, `rpeak-helper.go`).

## Features
- User management
- Device management
- Hospital management
- Medical record handling
- ECG signal processing utilities

## Prerequisites
- Go 1.20 or later
- A configured database (e.g., PostgreSQL, MySQL, etc.)

## Setup Instructions
1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd khidmaat-backend
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Configure the database connection in `config/db.go`.

4. Run the application:
   ```bash
   go run main.go
   ```

## API Endpoints
The application exposes RESTful APIs for various operations. Below are the main modules and their routes:

- **User Management**: `routers/user-routes.go`
- **Hospital Management**: `routers/hospital-routes.go`
- **ECG Data**: `routers/ecg-routes.go`

## Contributing
Contributions are welcome! Please follow these steps:
1. Fork the repository.
2. Create a new branch for your feature or bugfix.
3. Commit your changes and push the branch.
4. Submit a pull request.

## License
This project is licensed under the MIT License. See the LICENSE file for details.

## Contact
For any inquiries or support, please contact the project maintainers.
