FROM python:3.11-slim

# Set the working directory
WORKDIR /app

# Copy the requirements.txt file if you have additional dependencies
# COPY requirements.txt .

# Install Flask (and any other dependencies)
RUN pip install --no-cache-dir flask

# Copy the application code
COPY main.py .

# Set the environment variable for Flask
ENV FLASK_APP=main.py

# Expose the port the app runs on
EXPOSE 5000

# Command to run the application
CMD ["flask", "run", "--host=0.0.0.0", "--port=5000"]
