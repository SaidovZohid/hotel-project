package postgres

import (
	"fmt"
	"log"

	"github.com/SaidovZohid/hotel-project/storage/repo"
	"github.com/jmoiron/sqlx"
)

type hotelRepo struct {
	db *sqlx.DB
}

func NewHotelStorage(db *sqlx.DB) repo.HotelStorageI {
	return &hotelRepo{
		db: db,
	}
}

func (hd *hotelRepo) Create(h *repo.Hotel) (int64, error) {
	tr, err := hd.db.Begin()
	defer tr.Rollback()
	if err != nil {
		return 0, err
	}
	query := `
		INSERT INTO hotels (
			manager_id,
			hotel_name,
			description,
			address,
			image_url,
			num_of_rooms
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	err = tr.QueryRow(
		query,
		h.ManagerID,
		h.HotelName,
		h.Description,
		h.Address,
		h.ImageUrl,
		h.NumOfRooms,
	).Scan(&h.ID)
	if err != nil {
		return 0, err
	}

	queryImages := `
		INSERT INTO hotel_images (
			hotel_id,
			image_url,
			sequence_number
		) VALUES ($1, $2, $3)
	`
	for _, image := range h.Images {
		_, err = tr.Exec(
			queryImages,
			h.ID,
			image.ImageUrl,
			image.SequenceNumber,
		)
		if err != nil {
			return 0, err
		}
	}
	err = tr.Commit()
	if err != nil {
		return 0, err
	}
	return h.ID, nil
}

func (hd *hotelRepo) Get(hotel_id int64) (*repo.Hotel, error) {
	query := `
		SELECT 
			id,
			manager_id,
			hotel_name,
			description,
			address,
			image_url,
			num_of_rooms
		FROM hotels WHERE id = $1
	`
	var res repo.Hotel
	err := hd.db.QueryRow(query, hotel_id).Scan(
		&res.ID,
		&res.ManagerID,
		&res.HotelName,
		&res.Description,
		&res.Address,
		&res.ImageUrl,
		&res.NumOfRooms,
	)
	if err != nil {
		return nil, err
	}

	queryImages := `
		SELECT 
			id,
			hotel_id,
			image_url,
			sequence_number
		FROM hotel_images WHERE hotel_id = $1
	`
	rows, err := hd.db.Query(queryImages, res.ID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var image repo.HotelImage
		err := rows.Scan(
			&image.ID,
			&image.HotelID,
			&image.ImageUrl,
			&image.SequenceNumber,
		)
		if err != nil {
			return nil, err
		}

		res.Images = append(res.Images, &image)
	}

	return &res, nil
}

func (hd *hotelRepo) Update(h *repo.Hotel) error {
	tr, err := hd.db.Begin()
	defer tr.Rollback()
	if err != nil {
		return err
	}

	queryImageDelete := "DELETE FROM hotel_images WHERE hotel_id = $1"
	_, err = tr.Exec(queryImageDelete, h.ID)
	if err != nil {
		return err
	}

	query := `
		UPDATE hotels SET
			hotel_name = $1,
			description = $2,
			address = $3,
			image_url = $4,
			num_of_rooms = $5
		WHERE id = $6
	`
	_, err = tr.Exec(
		query,
		h.HotelName,
		h.Description,
		h.Address,
		h.ImageUrl,
		h.NumOfRooms,
		h.ID,
	)
	if err != nil {
		return err
	}

	queryImage := `
		INSERT INTO hotel_images (
			hotel_id,
			image_url,
			sequence_number
		) VALUES ($1, $2, $3)
	`
	for _, image := range h.Images {
		_, err = tr.Exec(
			queryImage,
			image.HotelID,
			image.ImageUrl,
			image.SequenceNumber,
		)
		if err != nil {
			return err
		}
	}
	log.Println("fffwf")
	err = tr.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (hd *hotelRepo) Delete(hotel_id int64) error {
	queryImageDelete := "DELETE FROM hotel_images WHERE hotel_id = $1"
	_, err := hd.db.Exec(queryImageDelete, hotel_id)
	if err != nil {
		return err
	}

	query := "DELETE FROM hotels WHERE id = $1"
	_, err = hd.db.Exec(query, hotel_id)
	if err != nil {
		return err
	}

	return nil
}

func (hd *hotelRepo) GetAll(params *repo.GetAllHotelsParams) (*repo.GetAllHotels, error) {
	offset := (params.Page - 1) * params.Limit 
	filter := " WHERE true "
	limit := fmt.Sprintf(" LIMIT %d OFFSET %d ", params.Limit, offset)
	if params.Search != "" {
		str := "%" + params.Search + "%"
		filter += fmt.Sprintf(" AND hotel_name ILIKE '%s' OR description ILIKE '%s' OR address ILIKE '%s' ", str, str, str)
	}
	if params.NumOfRooms != 0 {
		filter += fmt.Sprintf(" AND num_of_rooms = %d ", params.NumOfRooms)
	}

	var res repo.GetAllHotels
	res.Hotels = make([]*repo.Hotel, 0)
	queryImages := `
		SELECT 
			id,
			hotel_id,
			image_url,
			sequence_number
		FROM hotel_images
	`
	rowsImage, err := hd.db.Query(queryImages)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT 
			id,
			manager_id,
			hotel_name,
			description,
			address,
			image_url,
			num_of_rooms
		FROM hotels
	` + filter + limit
	rows, err := hd.db.Query(query)
	if err != nil {
		return nil, err	
	}
	for rows.Next() {
		var hotel repo.Hotel
		err = rows.Scan(
			&hotel.ID,
			&hotel.ManagerID,
			&hotel.HotelName,
			&hotel.Description,
			&hotel.Address,
			&hotel.ImageUrl,
			&hotel.NumOfRooms,
		)
		if err != nil {
			return nil, err
		}

	    for rowsImage.Next() {
			var image repo.HotelImage
			err = rowsImage.Scan(
				&image.ID,
				&image.HotelID,
				&image.ImageUrl,
				&image.SequenceNumber,
			)
			if err != nil {
				return nil, err
			}

			hotel.Images = append(hotel.Images, &image)
		}
		res.Hotels = append(res.Hotels, &hotel)
	}
	queryCount := " SELECT count(1) FROM hotels " + filter
	err = hd.db.QueryRow(queryCount).Scan(&res.Count)
	if err != nil {
		return nil, err
	}

	return &res, nil
}