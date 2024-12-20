# Remote Sensing Image Segmentation Web Application

This is a web application for semantic segmentation of remote sensing images. It consists of a React frontend and a Go backend.

## Project Structure
```
.
├── frontend/           # React frontend application
└── backend/           # Go backend server
```

## Setup Instructions

### Backend Setup
1. Navigate to the backend directory:
   ```bash
   cd backend
   ```
2. Run the Go server:
   ```bash
   go run main.go
   ```
   The server will start on port 8080.

### Frontend Setup
1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```
2. Install dependencies:
   ```bash
   npm install
   ```
3. Start the development server:
   ```bash
   npm start
   ```
   The frontend will be available at http://localhost:3000

## Usage
1. Open the web application in your browser
2. Click "Select Image" to choose a remote sensing image
3. Click "Process Image" to send it to the backend for segmentation
4. The segmented result will be displayed below the original image

## Note
This is a basic implementation. The current version includes:
- Image upload functionality
- Basic frontend UI with Material-UI components
- Backend API endpoint for file upload
- CORS support

To make this a production-ready application, you would need to:
1. Add proper error handling
2. Implement actual image segmentation using a machine learning model
3. Add input validation
4. Implement proper file storage
5. Add security measures
6. Add loading states and progress indicators
