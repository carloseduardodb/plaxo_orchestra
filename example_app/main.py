# Exemplo de aplicação FastAPI para testar Agent Spread
from fastapi import FastAPI

app = FastAPI(title="E-commerce API")

@app.get("/")
def read_root():
    return {"message": "E-commerce API"}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
