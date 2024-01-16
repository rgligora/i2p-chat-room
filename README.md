# I2P Encrypted Chat Room App

## Overview
This project is an eepsite-based chat room application utilizing Flask and Flask-SocketIO for secure and anonymous communication. It features end-to-end encryption using AES-GCM and RSA algorithms, facilitating secure chat rooms where users can exchange messages in real-time.

## Features
- Anonymous and secure chat rooms on the I2P network.
- End-to-end encryption with RSA and AES-GCM for message security.
- Real-time messaging through WebSockets.
- Configured to run as an eepsite with Nginx reverse proxy.

## Technology Stack
- **Backend:** Flask, Flask-SocketIO
- **Frontend:** HTML, CSS, JavaScript
- **Encryption:** Cryptography library in Python, Web Crypto API
- **Server:** Nginx (reverse proxy configuration for eepsite)


## Setup

### Prerequisites
- I2P Router installed and configured
- Python 3.x
- Nginx

### Installation

### Steps
1. **Clone the repository @ /var/www/ :**

   ```bash
   cd /var/www/
   git clone https://github.com/rgligora/i2p-chat-room-app.git

2. **Navigate to the project directory:**
    ```bash
   cd i2p-chat-room-app

3. **Venv**
   ```bash
   python3 -m venv venv
   source venv/bin/activate
   pip install -r requirements.txt

4. **Create a symbolic link to this configuration file in the sites-enabled directory:**
   New terminal tab
   ```bash
   cd /var/www/i2p-chat-room-app
   sudo ln -s /var/www/i2p-chat-room-app/chat-app.conf /etc/nginx/sites-enabled

5. **Restart Nginx to apply the changes:**
   ```bash
   sudo service nginx restart

6. **Install and start the I2P router form the official website**

7. **Setup the proxy to port 4444**
    In Firefox. Settings > Network Settings > Maunal Proxy Configuration > HTTP Proxy: 127.0.0.1 Port: 4444

8. **Open the I2P Router Concole's Hidden Services Manager**
   In the Global Tunnel Control submenu start the "Tunnel Wizard".
   Choose a Server Tunnel setup and click next. Choose the HTTP server tunnel type and click next. Name it "I2p Chat Room" and click next. Host is 0.0.0.0 and port is 443. Select the "Automatically start tunnel when router starts" and clock finish. Wait unit the status of your newly created tunnel is green to go to the next step.

9. **Start the Flask app:**
   ```bash
    gunicorn -k eventlet -w 1 -b 127.0.0.1:443 main:app

10. **Open the eepsite**
   Open the destination shown on your I2P Hidden Services menu next to the tunnel that you created in step 8. Example:
    http://nw7ruavzbpwqybf4fdidoyceetwk7rc357q3jevkfvfn7j6hknfa.b32.i2p/

![alt text](https://github.com/rgligora/i2p-chat-room-app/blob/main/showcase/home-desktop.png?raw=true)
![alt text](https://github.com/rgligora/i2p-chat-room-app/blob/main/showcase/chat-desktop.png?raw=true)

<div style="width: 50%;">
<table>
  <tr>
    <td><img src="https://github.com/rgligora/i2p-chat-room-app/blob/main/showcase/home-mobile.png?raw=true" alt="alt text"></td>
    <td><img src="https://github.com/rgligora/i2p-chat-room-app/blob/main/showcase/chat-mobile.png?raw=true" alt="alt text"></td>
  </tr>
</table>
</table>
</div>

![alt text](https://github.com/rgligora/i2p-chat-room-app/blob/main/showcase/chat-user-list-desktop.png?raw=true)