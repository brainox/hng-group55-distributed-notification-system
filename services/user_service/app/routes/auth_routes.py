from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy.future import select
from app.models.user import User
from app.schemas.auth_schema import RegisterRequest, LoginRequest, TokenResponse, UserOut
from app.core.security import hash_password, verify_password, create_access_token
from app.core.db import async_session
from app.utils.response import response

router = APIRouter(prefix="/v1/users", tags=["users"])

async def get_db():
    async with async_session() as session:
        yield session


@router.post("/register", response_model=dict)
async def register(payload: RegisterRequest, db: AsyncSession = Depends(get_db)):
    try:
        result = await db.execute(select(User).where(User.email == payload.email))
        existing = result.scalar_one_or_none()
        if existing:
            raise HTTPException(status_code=400, detail="Email already registered")

        new_user = User(
            full_name=payload.name,
            email=payload.email,
            password_hash=hash_password(payload.password),
            push_token=payload.push_token,
            preferences=payload.preferences.dict() if payload.preferences else {"email": True, "push": True}
        )
        db.add(new_user)
        await db.commit()
        await db.refresh(new_user)

        return response(True, data=UserOut.from_orm(new_user).dict(), message="User registered successfully")

    except Exception as e:
        await db.rollback()
        raise HTTPException(status_code=500, detail=f"Internal server error: {e}")


@router.post("/login", response_model=dict)
async def login(payload: LoginRequest, db: AsyncSession = Depends(get_db)):
    result = await db.execute(select(User).where(User.email == payload.email))
    user = result.scalar_one_or_none()

    if not user or not verify_password(payload.password, user.password_hash):
        raise HTTPException(status_code=status.HTTP_401_UNAUTHORIZED, detail="Invalid credentials")

    token = create_access_token(subject=str(user.user_id))
    return response(True, data=TokenResponse(access_token=token).dict(), message="Login successful")
