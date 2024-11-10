
# AndroidLibSocksLite

**AndroidLibSocksLite** is a highly optimized, Go-based SOCKS5 server core developed specifically for Android devices. This library empowers Android applications to host SOCKS5 proxy servers, enabling secure tunneling and efficient management of network traffic. The core has been meticulously crafted to support seamless SOCKS5 functionality, provided as an AAR (Android Archive) package for easy integration into Android projects.

## Key Features
- **Create Secure SOCKS5 Servers**: Launch SOCKS5 servers on specific ports with robust username/password authentication for added security.
- **Comprehensive Logging and Metrics**: Capture detailed access logs and gather essential metrics for real-time monitoring.
- **Integrated HTTP APIs**: Expose API endpoints to monitor server status, client activity, and connection logs with ease.

## Build requirements
* JDK
* Android SDK , NDK
* Go
* gomobile

## Building the AAR

To integrate this library into an Android application, first build it as an AAR file. Follow these steps to set up and build the AAR: 
1. **Install Go Mobile Tools**
   ```bash
   go install golang.org/x/mobile/cmd/gomobile@latest
   go install golang.org/x/mobile/cmd/gobind@latest
   ```

2. **Initialize and Prepare the Go Mobile Environment**
   ```bash
   go get
   gomobile init
   go mod tidy -v
   go get golang.org/x/mobile/bind
   ```

3. **Build the AAR**
   ```bash
   gomobile bind -v -androidapi 21 -ldflags='-s -w' -o libSocksLite.aar ./pkg/socks
   ```
   This command will generate an AAR package that can be seamlessly integrated into your Android project.

## Core Functions

AndroidLibSocksLite offers the following core functions for managing SOCKS5 servers and monitoring their statuses:

## CheckCoreVersion

The `CheckCoreVersion` function returns the current core version of the library. It can be used to verify which version of the core is being utilized in your application.

### `StartSocksServers(host: String, jsonData: String): Error`

- **Purpose**: Initializes and launches multiple SOCKS5 servers based on JSON input data. The servers, along with an HTTP server for monitoring, are bound to the specified `host` address.
- **Parameters**:
    - `host`: The IP address where the SOCKS5 and HTTP servers will run.
    - `jsonData`: JSON-formatted string containing user credentials and port configurations.
- **Return**: Returns an error if there’s any issue during setup; otherwise, servers start successfully.

### `Shutdown(): Error`

- **Purpose**: Gracefully terminates the HTTP server and all SOCKS5 servers, ensuring all resources are released and connections closed.
- **Return**: Returns an error if an issue occurs during shutdown.

## HTTP API Endpoints

AndroidLibSocksLite also provides HTTP API endpoints to facilitate server monitoring, metrics retrieval, and log access.

### `/getClientStatus`
- **Method**: `GET`
- **Description**: Returns the current status of each user's SOCKS5 server, including port and running status.
- **Response Example**:
  ```json
  [
    {
      "username": "user1",
      "port": 8000,
      "running": true
    },
    ...
  ]
  ```

### `/shutdown`
- **Method**: `GET`
- **Description**: Shuts down all active SOCKS5 servers gracefully and terminates the HTTP server.
- **Response**: Plain text message indicating successful shutdown.

## Required Android Permissions

Add the following permissions to your Android `AndroidManifest.xml` to allow network access, which is essential for creating SOCKS5 connections:

```xml
<uses-permission android:name="android.permission.INTERNET" />
<uses-permission android:name="android.permission.ACCESS_NETWORK_STATE" />
```

## Android Integration Example

Below are examples of using AndroidLibSocksLite in an Android app, with Kotlin and Java code that utilizes Google Gson for JSON parsing.

### Kotlin Example

Define the `User` data class:

```kotlin
data class User(
    val username: String,
    val password: String,
    val port: Int
)
```

In your main activity:

```kotlin
import AndroidLibSocksLite.AndroidLibSocksLite
import com.google.gson.Gson

class MainActivity : AppCompatActivity() {

    private val core = AndroidLibSocksLite()

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        val users = listOf(
            User("user1", "pass1", 8000),
            User("user2", "pass2", 8001)
        )

        val jsonData = Gson().toJson(users)
        core.startSocksServers("127.0.0.1", jsonData)
    }

    override fun onDestroy() {
        super.onDestroy()
        core.shutdown() // Graceful shutdown
    }
}
```

### Java Example

Define the `User` class:

```java
public class User {
    private String username;
    private String password;
    private int port;

    public User(String username, String password, int port) {
        this.username = username;
        this.password = password;
        this.port = port;
    }

    // Getters and Setters
}
```

In your main activity:

```java
import AndroidLibSocksLite.AndroidLibSocksLite;
import com.google.gson.Gson;
import java.util.ArrayList;
import java.util.List;

public class MainActivity extends AppCompatActivity {

    private AndroidLibSocksLite core = new AndroidLibSocksLite();

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);

        List<User> users = new ArrayList<>();
        users.add(new User("user1", "pass1", 8000));
        users.add(new User("user2", "pass2", 8001));

        String jsonData = new Gson().toJson(users);
        core.startSocksServers("127.0.0.1", jsonData);
    }

    @Override
    protected void onDestroy() {
        super.onDestroy();
        core.shutdown(); // Graceful shutdown
    }
}
```

## Local Testing on Windows

For development or testing on a Windows machine, you can run a `main.go` file locally to start the servers.

### Sample `main.go`

```go
package main

import (
	"AndroidLibSocksLite"
	"fmt"
)

func main() {
	jsonData := `[
        {"username": "user1", "password": "pass1", "port": 8000},
        {"username": "user2", "password": "pass2", "port": 8001}
    ]`

	err := AndroidLibSocksLite.StartSocksServers("127.0.0.1", jsonData)
	if err != nil {
		fmt.Printf("Error starting servers: %v\n", err)
	} else {
		fmt.Println("Servers started successfully on localhost")
	}
}
```

To run this on Windows:
1. Save the code above as `main.go`.
2. Execute it with:
   ```bash
   go run main.go
   ```

This will start the SOCKS5 and HTTP servers locally, allowing you to test the APIs at `http://127.0.0.1:8080/getClientStatus` using a browser or HTTP client.

---

## Credits

This project utilizes the [go-socks5 library](https://github.com/armon/go-socks5) by Armon Dadgar.

## License

This project is licensed under the [Apache License 2.0](LICENSE).

---

Made with ❤️ by Tamim
