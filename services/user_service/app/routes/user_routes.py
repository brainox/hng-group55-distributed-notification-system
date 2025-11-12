from fastapi import APIRouter, Depends, HTTPException, Header, status
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import select
from uuid import UUID
from app.core.db import async_session
from app.models.user import User
from app.schemas.user_schema import UserCreate, UserOut, UserUpdate, UserPreference
from app.core.security import hash_password
from app.utils.response import response

router = APIRouter(prefix="/v1/users", tags=["users"])


async def get_db():
    async with async_session() as session:
        yield session


# ✅ Create a new user
@router.post("", response_model=dict)
async def create_user(
    payload: UserCreate,
    idempotency_key: str | None = Header(None, alias="idempotency_key"),
    db: AsyncSession = Depends(get_db)
):
    stmt = select(User).where(User.email == payload.email)
    existing = await db.scalar(stmt)
    if existing:
        raise HTTPException(status_code=400, detail="Email already exists")

    user = User(
        email=payload.email,
        full_name=payload.full_name,
        password_hash=hash_password(payload.password),
        preferences={"email": True, "push": True},
    )
    db.add(user)
    await db.commit()
    await db.refresh(user)
    return response(True, data=UserOut.from_orm(user), message="User created successfully")


# ✅ List all users
@router.get("", response_model=dict)
async def list_users(db: AsyncSession = Depends(get_db)):
    result = await db.execute(select(User))
    users = result.scalars().all()
    return response(True, data=[UserOut.from_orm(u) for u in users], message="Users retrieved successfully")


# ✅ Get single user
@router.get("/{user_id}", response_model=dict)
async def get_user(user_id: UUID, db: AsyncSession = Depends(get_db)):
    user = await db.get(User, user_id)
    if not user:
        raise HTTPException(status_code=404, detail="User not found")
    return response(True, data=UserOut.from_orm(user), message="User retrieved successfully")

# ✅ Update user (name, email, push_token, preferences)
@router.patch("/{user_id}", response_model=dict)
async def update_user(user_id: UUID, payload: UserUpdate, db: AsyncSession = Depends(get_db)):
    user = await db.get(User, user_id)
    if not user:
        raise HTTPException(status_code=404, detail="User not found")

    update_data = payload.model_dump(exclude_unset=True)

    for field, value in update_data.items():
        if field == "preferences" and value is not None:
            # handle both dict or Pydantic object safely
            new_prefs = value if isinstance(value, dict) else value.model_dump()
            # ensure preferences exist before updating
            user.preferences = {**(user.preferences or {}), **new_prefs}
        else:
            setattr(user, field, value)

    await db.commit()
    await db.refresh(user)
    return response(True, data=UserOut.from_orm(user), message="User updated successfully")



# ✅ Get user preferences
@router.get("/{user_id}/preferences", response_model=dict)
async def get_user_preferences(user_id: UUID, db: AsyncSession = Depends(get_db)):
    user = await db.get(User, user_id)
    if not user:
        raise HTTPException(status_code=404, detail="User not found")

    prefs = user.preferences or {"email": True, "push": True}
    return response(True, data=prefs, message="User preferences retrieved")


# ✅ Update user preferences
@router.put("/{user_id}/preferences", response_model=dict)
async def update_user_preferences(user_id: UUID, payload: UserPreference, db: AsyncSession = Depends(get_db)):
    user = await db.get(User, user_id)
    if not user:
        raise HTTPException(status_code=404, detail="User not found")

    user.preferences.update(payload.model_dump())
    await db.commit()
    await db.refresh(user)
    return response(True, data=user.preferences, message="User preferences updated successfully")
