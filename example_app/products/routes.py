from fastapi import APIRouter
from .models import Product

router = APIRouter(prefix="/products")

@router.get("/")
def list_products():
    return {"products": []}

@router.get("/{product_id}")
def get_product(product_id: int):
    return {"product": {"id": product_id}}

@router.post("/")
def create_product(name: str, price: float):
    return {"message": "Product created"}
