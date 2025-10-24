import { UserResponseDto } from './user-response.dto';
import { CreateUserDto, UpdateUserDto } from './user.dto';

export abstract class UserServiceAbstract {
  constructor() {}

  abstract findAll(): Promise<UserResponseDto[]>;
  abstract findOne(id: string): Promise<UserResponseDto | null>;
  abstract create(createUserDto: CreateUserDto): Promise<UserResponseDto>;
  abstract update(
    id: string,
    updateUserDto: UpdateUserDto,
  ): Promise<UserResponseDto | null>;
  abstract remove(id: string): Promise<void>;
}
