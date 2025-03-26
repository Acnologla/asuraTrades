# Trade System

## Overview

This project is a trade system that allows users to exchange items securely and efficiently. It consists of two main services: the **Token Service** and the **Trade Service**. The system follows a **hexagonal architecture** and uses **PostgreSQL** for data storage, a **local cache** for performance optimization, and **JWT** for token generation.

## Services

### Token Service

- Generates tokens using two user IDs.
- Retrieves user information based on a given token.
- Authenticates trades and stores trade IDs.

### Trade Service

- Facilitates item exchanges between users.
- Validates whether users own the items and if the items are tradable.
- Maintains the state of a trade.
- Allows users to create new trades.
- Enables trade confirmation, providing a **5-second window** during which users can cancel the trade if permitted.

### gRPC Trusty Trades

- Implements a secure trading mechanism using gRPC.
- Uses a unique token for authentication.
- Sends all trade items simultaneously for trading.
- Provides an all-or-nothing trade process.

## Technology Stack

- **Database:** PostgreSQL
- **Cache:** Local cache
- **Authentication:** JWT
- **Communication Protocols:**
  - **HTTP Controllers**: Used for token creation and retrieving user information via tokens.
  - **WebSockets**: Used for real-time trade interactions, including adding/removing items and confirming trades.
  - **gRPC**: Implements trades where all items are sent at once, without the need for confirmation.

## Architecture

This project follows the **hexagonal architecture** pattern to ensure modularity and scalability. The main components include:

- **Adapters:** they handle external dependencies such as databases, caching, WebSockets, and gRPC.
- **Core Domain:** this contains the business logic and domain models.
- **Entry Point:** The application starts from `cmd/asura-trades/main.go`.

## Testing

The project includes comprehensive unit tests where we mock the ports (interfaces) using gomock. This approach enables us to thoroughly test our business logic in isolation from external dependencies, ensuring the reliability of our services.

## License

This project is licensed under the MIT License.

### Notes

- Consider implementing a more robust caching layer if scaling is required.
