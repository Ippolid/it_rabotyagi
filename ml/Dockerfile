FROM python:3.11-slim

WORKDIR /app

# Install dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy application code
COPY interview_trainer.py .
COPY voprosy.last.csv .

# Expose port
EXPOSE 8001

# Run the application
CMD ["uvicorn", "interview_trainer:app", "--host", "0.0.0.0", "--port", "8001"]
