from fastapi import APIRouter, Depends
from .models import User

router = APIRouter(prefix="/auth")

@router.post("/login")
def login(email: str, password: str):
    return {"token": "fake-jwt-token"}

@router.post("/register")
def register(email: str, password: str):
    return {"message": "User created"}
