

# E-Commerce Project README.md

## Overview
This project is a Go-based e-commerce platform that handles product management, order processing, and user interactions. It leverages PostgreSQL for database operations using GORM as the ORM (Object-Relational Mapping) tool.

## Project Structure
The project's structure is organized into several key packages:
- `storage`: Contains the generic CRUD (Create, Read, Update, Delete) operations.
- `types`: Defines data structures and models used throughout the application.
- `user`: Implements user-related functionalities including authentication and profile management.
- `product` & `order`: Handle product listings, inventory management, order creation, and tracking.

## Dependencies
The project relies on the following key dependencies:
- **GORM**: For database interactions.
- **PostgreSQL**: As the primary database backend.

## Getting Started

### Prerequisites
Ensure you have Go installed. PostgreSQL is also required for the application to function correctly.

### Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/ecomm.git
   cd ecomm
   ```
2. Install dependencies using Go modules:
   ```bash
   go mod tidy
   ```

### Setup
1. **Database Configuration**:
   - Create a PostgreSQL database.
   - Update the connection details in your `.env` file (if used) or directly in the application configuration.

2. **Running Migrations**:
   The SQL schema for orders and products is included. You can run these migrations using a tool like `goose` or manually execute them against your database.

## Database Schema

### Orders Table
| Column Name       | Data Type          | Description                     |
|-------------------|--------------------|---------------------------------|
| id                | UUID               | Primary key, auto-generated     |
| user_id           | UUID               | Foreign key referencing users  |
| total             | FLOAT              | Total amount of the order      |
| status            | VARCHAR(50)        | Order status (e.g., pending, shipped) |
| address           | TEXT               | Shipping address               |
| created_at        | TIMESTAMP          | Record creation timestamp       |
| updated_at        | TIMESTAMP          | Last update timestamp          |

### Order Items Table
| Column Name       | Data Type          | Description                     |
|-------------------|--------------------|---------------------------------|
| id                | UUID               | Primary key, auto-generated     |
| order_id          | UUID               | Foreign key referencing orders  |
| product_id        | UUID               | Foreign key referencing products|
| quantity          | INT                | Number of items                 |
| price             | FLOAT              | Price per item                  |
| created_at        | TIMESTAMP          | Record creation timestamp       |
| updated_at        | TIMESTAMP          | Last update timestamp          |

## Contributing
Contributions are welcome! For major changes, please open an issue first to discuss what you'd like to change. Please ensure your code adheres to Go's best practices and includes tests.

## License
This project is licensed under the MIT License - see the `LICENSE.md` file for details.

---

For more detailed information about specific functionalities or further documentation, check out the respective package directories.