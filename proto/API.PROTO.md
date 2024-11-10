### Telegram Mini App API:

```proto
syntax = "proto3";

package telegram_mini_app;

// Message representing a user in the Telegram Mini App
message User {
    int64 user_id = 1;
    string username = 2;
    string first_name = 3;
    string last_name = 4;
    string language_code = 5;
    bool is_bot = 6; // Whether the user is a bot or not
}

// Represents a simple text message
message Message {
    int64 message_id = 1;
    int64 chat_id = 2;
    int64 sender_id = 3;
    string text = 4;
    string timestamp = 5;
}

// Represents a chat
message Chat {
    int64 chat_id = 1;
    string chat_name = 2;
    string chat_type = 3; // "private", "group", "supergroup", etc.
    repeated User participants = 4;
}

// Response to send a message
message SendMessageResponse {
    bool success = 1;
    string error_message = 2;
    Message message = 3;
}

// Request to send a text message
message SendMessageRequest {
    int64 chat_id = 1;
    int64 sender_id = 2;
    string text = 3;
}

// Request to get user information
message GetUserInfoRequest {
    int64 user_id = 1;
}

// Response with user information
message GetUserInfoResponse {
    User user = 1;
}

// A request to create a new chat
message CreateChatRequest {
    string chat_name = 1;
    string chat_type = 2; // e.g., "private", "group"
    repeated int64 user_ids = 3;
}

// Response when a chat is created
message CreateChatResponse {
    int64 chat_id = 1;
    string chat_name = 2;
    string chat_type = 3;
    bool success = 4;
    string error_message = 5;
}

// Request for app initialization with init_data
message InitDataRequest {
    string init_data = 1;
}

// Response for app initialization
message InitDataResponse {
    User user = 1;  // User information parsed from init_data
    Chat chat = 2;  // Chat information parsed from init_data
    string bot_name = 3; // Bot information, if applicable
    string payload = 4;  // Optional custom payload
    bool success = 5;
    string error_message = 6;
}

// Response with a list of all chats the user is part of
message GetChatsResponse {
    repeated Chat chats = 1;
}

// Response with a list of messages from a particular chat
message GetMessagesResponse {
    repeated Message messages = 1;
}

// Request to send a media message (e.g., photo, video)
message SendMediaMessageRequest {
    int64 chat_id = 1;
    int64 sender_id = 2;
    string media_url = 3;  // URL of the media file
    string media_type = 4;  // Type of media (e.g., "photo", "video")
    string caption = 5;  // Optional caption
}

// Response for sending media message
message SendMediaMessageResponse {
    bool success = 1;
    string error_message = 2;
    Message message = 3;
}

// Request to delete a message
message DeleteMessageRequest {
    int64 chat_id = 1;
    int64 message_id = 2;
}

// Response to delete a message
message DeleteMessageResponse {
    bool success = 1;
    string error_message = 2;
}

// Request to get the bot's information
message GetBotInfoRequest {}

message BotInfo {
    string bot_name = 1;
    string bot_username = 2;
    string bot_description = 3;
}

message GetBotInfoResponse {
    BotInfo bot_info = 1;
    bool success = 2;
    string error_message = 3;
}

// Request to join a chat
message JoinChatRequest {
    int64 chat_id = 1;
    int64 user_id = 2;
}

// Response for joining a chat
message JoinChatResponse {
    bool success = 1;
    string error_message = 2;
}

// Request to leave a chat
message LeaveChatRequest {
    int64 chat_id = 1;
    int64 user_id = 2;
}

// Response for leaving a chat
message LeaveChatResponse {
    bool success = 1;
    string error_message = 2;
}

// Request to pin a message
message PinMessageRequest {
    int64 chat_id = 1;
    int64 message_id = 2;
}

// Response for pinning a message
message PinMessageResponse {
    bool success = 1;
    string error_message = 2;
}

// Request to unpin a message
message UnpinMessageRequest {
    int64 chat_id = 1;
    int64 message_id = 2;
}

// Response for unpinning a message
message UnpinMessageResponse {
    bool success = 1;
    string error_message = 2;
}

// Payment-related messages (for in-app purchases or donations)
message PaymentInfo {
    string payment_id = 1;
    double amount = 2;  // Payment amount
    string currency = 3;  // Currency, e.g., "USD"
    string description = 4;  // Description of the payment (item, service, etc.)
}

message ProcessPaymentRequest {
    string user_id = 1;
    PaymentInfo payment_info = 2;
}

message ProcessPaymentResponse {
    bool success = 1;
    string error_message = 2;
    string payment_url = 3;  // If payment URL is generated
}

// Telegram Mini App Service Definition
service TelegramMiniApp {

    // Method to initialize the Mini App with init_data
    rpc InitData(InitDataRequest) returns (InitDataResponse);

    // Method to send a text message
    rpc SendMessage(SendMessageRequest) returns (SendMessageResponse);

    // Method to get user info by user ID
    rpc GetUserInfo(GetUserInfoRequest) returns (GetUserInfoResponse);

    // Method to create a new chat (private or group)
    rpc CreateChat(CreateChatRequest) returns (CreateChatResponse);

    // Method to get all chats the user is part of
    rpc GetChatsForUser(int64) returns (GetChatsResponse);

    // Method to get messages from a chat
    rpc GetMessagesFromChat(int64) returns (GetMessagesResponse);

    // Method to send media messages (photo, video, etc.)
    rpc SendMediaMessage(SendMediaMessageRequest) returns (SendMediaMessageResponse);

    // Method to delete a message
    rpc DeleteMessage(DeleteMessageRequest) returns (DeleteMessageResponse);

    // Method to get bot's information
    rpc GetBotInfo(GetBotInfoRequest) returns (GetBotInfoResponse);

    // Method to join a chat
    rpc JoinChat(JoinChatRequest) returns (JoinChatResponse);

    // Method to leave a chat
    rpc LeaveChat(LeaveChatRequest) returns (LeaveChatResponse);

    // Method to pin a message
    rpc PinMessage(PinMessageRequest) returns (PinMessageResponse);

    // Method to unpin a message
    rpc UnpinMessage(UnpinMessageRequest) returns (UnpinMessageResponse);

    // Method to process a payment
    rpc ProcessPayment(ProcessPaymentRequest) returns (ProcessPaymentResponse);
}
```

### Functions:

1. **SendMediaMessage**:

   - Allows sending media (like photos, videos, etc.) to a chat.
   - Useful for sending rich content beyond just text messages.

2. **DeleteMessage**:
   - Allows deleting a specific message from a chat. This could be useful for chat moderators or the app itself to remove content.
3. **GetBotInfo**:

   - Retrieves information about the bot itself (name, username, description). This can help the Mini App present information about the bot or handle certain actions.

4. **JoinChat & LeaveChat**:
   - Methods to allow users to join or leave a specific chat. This can be useful for managing group chats or handling specific actions like leaving a support group or joining a channel.
5. **PinMessage & UnpinMessage**:

   - Methods for pinning or unpinning messages in a group chat. Useful for managing important announcements or messages within the chat.

6. **ProcessPayment**:
   - Allows the Mini App to process payments, such as donations or in-app purchases.
   - This could include handling payment gateways, generating URLs for payments, or processing financial transactions.

### Why These Functions Are Useful:

- **Media and Rich Content**: The ability to send media messages (images, videos) adds a richer user experience, which is often expected in modern applications.
- **Message Management**: Functions like deleting, pinning, and unpinning messages are vital for chat moderation, especially in larger groups or channels.
- **Payments**: Handling payments directly within a Mini App can enable features like in-app purchases or donations, making it useful for e-commerce, content creators, or fundraising.
- **Chat Management**: Joining or leaving chats is useful for user engagement, especially in dynamic environments like community-driven groups or events.
- **Bot Info**: Retrieving bot details helps personalize the user experience, making interactions with the bot more informative and seamless.

### Final Thoughts:

This extended `.proto` file now covers a wide range of potential functionality needed in a **Telegram Mini App**, from simple messaging to media handling, user management, and even payments. Depending on your Mini App's specific requirements, you can further extend or modify this API to suit your needs.
