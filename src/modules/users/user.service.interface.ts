import { User } from './user.entity';
import { CreateUserDto, UpdateUserDto } from './user.dto';

export abstract class UserServiceAbstract {
  constructor() {}

  abstract findAll(): Promise<User[]>;
  abstract findOne(id: string): Promise<User | null>;
  abstract create(createUserDto: CreateUserDto): Promise<User>;
  abstract update(
    id: string,
    updateUserDto: UpdateUserDto,
  ): Promise<User | null>;
  abstract remove(id: string): Promise<void>;
}
