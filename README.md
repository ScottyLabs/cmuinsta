# cmuinsta
Automating frontend for the cmu prefrosh prage

### Getting Started
To get started, follow these steps:

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/cmuinsta.git
   ```
2. Make sure you have NixOS with DirEnv installed, similar if not same as [Terrier contribution setup](https://github.com/ScottyLabs/terrier/blob/main/CONTRIBUTING.md)
3. Navigate to the project directory:
   ```bash
   cd cmuinsta
   ```
4. Copy the `.env.sample` file and modify
5. Initialize the database
   ```bash
   just db-init
   ```
6. Run the app for development or production
   ```bash
   just dev
   just up
   ```
7. Kill the app
   ```bash
   just down
   ```
