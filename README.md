# cmuinsta
Automating frontend for the cmu prefrosh prage

### Getting Started
To get started, follow these steps:

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/cmuinsta.git
   ```
2. Make sure you have NixOS with DirEnv installed 
3. Navigate to the project directory:
   ```bash
   cd cmuinsta
   ```
4. Initialize the database
   ```bash
   just init-db
   ```
5. Run the app for development or production
   ```bash
   just dev
   just up
   ```
6. Kill the app
   ```bash
   just down
   ```
