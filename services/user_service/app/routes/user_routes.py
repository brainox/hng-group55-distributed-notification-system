from fastapi import APIRouter, Depends, HTTPException, Header, status
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import select
from uuid import UUID
from app.core.db import async_session
from app.models.user import User
from app.schemas.user_schema import UserCreate, UserOut
from app.core.security import hash_password
from app.utils.response import response

router = APIRouter(prefix="/v1/users", tags=["users"])

async def get_db():
    async with async_session() as session:
        yield session


@router.post("", response_model=dict)
async def create_user(
    payload: UserCreate,
    idempotency_key: str | None = Header(None, alias="idempotency_key"),
    db: AsyncSession = Depends(get_db)
):
    stmt = select(User).where(User.email == payload.email)
    existing = await db.scalar(stmt)
    if existing:
        raise HTTPException(status_code=400, detail="email_exists")

    user = User(
        email=payload.email,
        full_name=payload.full_name,
        password_hash=hash_password(payload.password)
    )
    db.add(user)
    await db.commit()
    await db.refresh(user)
    return response(True, data=UserOut.from_orm(user), message="User created successfully")


@router.get("/{user_id}", response_model=dict)
async def get_user(user_id: UUID, db: AsyncSession = Depends(get_db)):
    user = await db.get(User, user_id)
    if not user:
        raise HTTPException(status_code=404, detail="User not found")
    return response(True, data=UserOut.from_orm(user), message="User retrieved successfully")
