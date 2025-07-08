package persistance

import "taskmgmtsystem/internal/core/session"

type SessionRepo struct {
	db *Database
}

func NewSessionRepo(d *Database) session.SessionRepoImpl {
	return SessionRepo{db: d}
}

func (u SessionRepo) CreateSession(session session.Session) error {
	_, err := u.db.db.Exec(`
		INSERT INTO SESSIONS (ID, USERID, TOKENHASH, ISSUED_AT,EXPIRES_AT)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (USERID) DO UPDATE 
		SET ID = EXCLUDED.ID, 
		    TOKENHASH = EXCLUDED.TOKENHASH,
		    EXPIRES_AT = EXCLUDED.EXPIRES_AT,
		    ISSUED_AT = EXCLUDED.ISSUED_AT;
	`, session.Id, session.UserId, session.TokenHash, session.IssuedAt, session.ExpiresAt)
	if err != nil {
		return err
	}
	return nil
}

func (u SessionRepo) DeleteSession(userId int) error {
	query := "DELETE FROM SESSIONS WHERE userid=$1"

	_, err := u.db.db.Query(query, userId)
	if err != nil {
		return err
	}

	return nil
}
