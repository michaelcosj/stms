package repository

const (
	insertUserStmt = `
    INSERT INTO users 
    (email, username, password, time_created) 
    VALUES (?, ?, ?, ?)
  `
	selectUserByIDStmt = `
    SELECT user_id, email, username, password, is_verified
    FROM users WHERE user_id = ?
  `

	selectUserByEmailStmt = `
    SELECT user_id, email, username, password, is_verified
    FROM users WHERE email = ?
  `

	updateUserStmt = `
    UPDATE users SET username = ?, is_verified = ?
    WHERE user_id = ?
  `

	deleteUserStmt = `
    DELETE FROM users
    WHERE user_id = ?
  `

	insertTaskStmt = `
    INSERT INTO tasks
    (name, tag, priority, is_completed, description, time_due,
      time_created, user_id)
    VALUES (?, ?, ?, ?, ?, ?, ?, ?)
  `

	selectTasksStmt = `
    SELECT task_id, name, tag, priority, is_completed, description,
      time_due, time_created, time_completed
    FROM tasks WHERE user_id = ?
  `

	updateTaskStmt = `
    UPDATE tasks SET name = ?, tag = ?, priority = ?, is_completed = ?,
      description = ?, time_due = ?
    WHERE task_id = ? 
  `

	deleteTaskStmt = `
    DELETE FROM tasks 
    WHERE task_id = ?
  `
)
