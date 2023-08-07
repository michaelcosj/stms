package repository

const (
	insertUserCommand = `
    INSERT INTO users 
    (email, username, password, time_created) 
    values (?, ?, ?, ?)
  `
	selectUserByIDCommand = `
    SELECT id, email, username, password, is_verified
    FROM users 
    WHERE id = ?
  `

	selectUserByEmailCommand = `
    SELECT id, email, username, password, is_verified
    FROM users 
    WHERE email = ?
  `

	selectUserTasksCommand = `
    SELECT * 
    FROM tasks
    WHERE user_id = ?
  `

	updateUserCommand = `
    UPDATE users SET 
    username = ?, is_verified = ?
    WHERE id = ?
  `
)
