from sqlalchemy import Column, Integer, String, Float, Text

class Product:
    __tablename__ = "products"
    
    id = Column(Integer, primary_key=True)
    name = Column(String(100))
    description = Column(Text)
    price = Column(Float)
    category_id = Column(Integer)
