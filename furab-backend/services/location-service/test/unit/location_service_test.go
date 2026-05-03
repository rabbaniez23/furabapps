package unit

import (
	"context"
	"errors"
	"testing"
	"time"

	"furab-backend/services/location-service/internal/model"
	"furab-backend/services/location-service/internal/repository/mock_repository"
	"furab-backend/services/location-service/internal/service"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUpdateLocation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockLocationRepository(ctrl)
	svc := service.NewLocationService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		req := model.UpdateLocationRequest{
			DriverID:  "driver1",
			Latitude:  -6.200000,
			Longitude: 106.816666,
			Timestamp: time.Now(),
		}

		mockRepo.EXPECT().UpdateLocation(gomock.Any(), req).Return(nil)

		err := svc.UpdateLocation(context.Background(), req)
		assert.NoError(t, err)
	})

	t.Run("Error - Koordinat tidak valid", func(t *testing.T) {
		req := model.UpdateLocationRequest{
			DriverID:  "driver1",
			Latitude:  200,
			Longitude: 300,
			Timestamp: time.Now(),
		}

		// repo MUST NOT be called
		mockRepo.EXPECT().UpdateLocation(gomock.Any(), gomock.Any()).Times(0)

		err := svc.UpdateLocation(context.Background(), req)
		assert.Error(t, err)
		assert.Equal(t, "invalid coordinate", err.Error())
	})
}

func TestUpdateStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockLocationRepository(ctrl)
	svc := service.NewLocationService(mockRepo)

	t.Run("Success - Available", func(t *testing.T) {
		req := model.UpdateStatusRequest{
			DriverID:     "driver1",
			DriverStatus: "available",
		}

		mockRepo.EXPECT().UpdateStatus(gomock.Any(), req).Return(nil)

		err := svc.UpdateStatus(context.Background(), req)
		assert.NoError(t, err)
	})

	t.Run("Success - Busy", func(t *testing.T) {
		req := model.UpdateStatusRequest{
			DriverID:     "driver1",
			DriverStatus: "busy",
		}

		mockRepo.EXPECT().UpdateStatus(gomock.Any(), req).Return(nil)

		err := svc.UpdateStatus(context.Background(), req)
		assert.NoError(t, err)
	})

	t.Run("Error - Status tidak valid", func(t *testing.T) {
		req := model.UpdateStatusRequest{
			DriverID:     "driver1",
			DriverStatus: "offline",
		}

		mockRepo.EXPECT().UpdateStatus(gomock.Any(), gomock.Any()).Times(0)

		err := svc.UpdateStatus(context.Background(), req)
		assert.Error(t, err)
		assert.Equal(t, "invalid status", err.Error())
	})
}

func TestSearchNearbyDrivers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockLocationRepository(ctrl)
	svc := service.NewLocationService(mockRepo)

	t.Run("Success - Driver ditemukan", func(t *testing.T) {
		req := model.SearchDriverRequest{
			LatitudeOrigin:  -6.200000,
			LongitudeOrigin: 106.816666,
			Radius:          5,
		}

		mockGeos := []redis.GeoLocation{
			{Name: "driver1", Longitude: 106.816666, Latitude: -6.200000, Dist: 1.2},
		}

		mockRepo.EXPECT().SearchNearbyDrivers(gomock.Any(), req).Return(mockGeos, nil)
		mockRepo.EXPECT().IsDriverActive(gomock.Any(), "driver1").Return(true, nil)
		mockRepo.EXPECT().GetStatus(gomock.Any(), "driver1").Return("available", nil)

		drivers, err := svc.SearchNearbyDrivers(context.Background(), req)
		assert.NoError(t, err)
		assert.Len(t, drivers, 1)
		assert.Equal(t, "driver1", drivers[0].DriverID)
		assert.Equal(t, "available", drivers[0].DriverStatus)
	})

	t.Run("Success - Tidak ada driver", func(t *testing.T) {
		req := model.SearchDriverRequest{
			LatitudeOrigin:  -6.200000,
			LongitudeOrigin: 106.816666,
			Radius:          1,
		}

		mockRepo.EXPECT().SearchNearbyDrivers(gomock.Any(), req).Return([]redis.GeoLocation{}, nil)

		drivers, err := svc.SearchNearbyDrivers(context.Background(), req)
		assert.NoError(t, err)
		assert.Len(t, drivers, 0)
	})

	t.Run("Filtering Busy Driver", func(t *testing.T) {
		req := model.SearchDriverRequest{
			LatitudeOrigin:  -6.200000,
			LongitudeOrigin: 106.816666,
			Radius:          5,
		}

		mockGeos := []redis.GeoLocation{
			{Name: "driver1", Longitude: 106.816666, Latitude: -6.200000, Dist: 1.2},
			{Name: "driver2", Longitude: 106.816666, Latitude: -6.200000, Dist: 1.5},
		}

		mockRepo.EXPECT().SearchNearbyDrivers(gomock.Any(), req).Return(mockGeos, nil)
		
		mockRepo.EXPECT().IsDriverActive(gomock.Any(), "driver1").Return(true, nil)
		mockRepo.EXPECT().GetStatus(gomock.Any(), "driver1").Return("busy", nil) // Should be filtered
		
		mockRepo.EXPECT().IsDriverActive(gomock.Any(), "driver2").Return(true, nil)
		mockRepo.EXPECT().GetStatus(gomock.Any(), "driver2").Return("available", nil) // Should pass

		drivers, err := svc.SearchNearbyDrivers(context.Background(), req)
		assert.NoError(t, err)
		assert.Len(t, drivers, 1)
		assert.Equal(t, "driver2", drivers[0].DriverID)
	})

	t.Run("Error - Input tidak valid", func(t *testing.T) {
		req := model.SearchDriverRequest{
			LatitudeOrigin:  999, // Invalid
			LongitudeOrigin: 999, // Invalid
			Radius:          5,
		}

		mockRepo.EXPECT().SearchNearbyDrivers(gomock.Any(), gomock.Any()).Times(0)

		drivers, err := svc.SearchNearbyDrivers(context.Background(), req)
		assert.Error(t, err)
		assert.Equal(t, "invalid input", err.Error())
		assert.Nil(t, drivers)
	})
}

func TestTrackDriver(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockLocationRepository(ctrl)
	svc := service.NewLocationService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		driverID := "driver1"
		ts := time.Now()
		mockResponse := &model.TrackLocationResponse{
			DriverID:  driverID,
			Latitude:  -6.200000,
			Longitude: 106.816666,
			Timestamp: ts,
		}

		mockRepo.EXPECT().TrackDriver(gomock.Any(), driverID).Return(mockResponse, nil)

		res, err := svc.TrackDriver(context.Background(), driverID)
		assert.NoError(t, err)
		assert.Equal(t, mockResponse, res)
	})

	t.Run("Error - Driver tidak ditemukan", func(t *testing.T) {
		driverID := "driver999"

		mockRepo.EXPECT().TrackDriver(gomock.Any(), driverID).Return(nil, errors.New("driver not found"))

		res, err := svc.TrackDriver(context.Background(), driverID)
		assert.Error(t, err)
		assert.Equal(t, "driver not found", err.Error())
		assert.Nil(t, res)
	})
}

func TestRequestEmergencyLocation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockLocationRepository(ctrl)
	svc := service.NewLocationService(mockRepo)

	t.Run("Success Case - emergency location ditemukan", func(t *testing.T) {
		driverID := "driver-emergency-1"
		ts := time.Now()

		// Equivalent call in current service: emergency location uses TrackDriver lookup.
		mockLocation := &model.TrackLocationResponse{
			DriverID:  driverID,
			Latitude:  -6.175392,
			Longitude: 106.827153,
			Timestamp: ts,
		}

		mockRepo.EXPECT().
			TrackDriver(gomock.Any(), gomock.Eq(driverID)).
			Return(mockLocation, nil).
			Times(1)

		res, err := svc.TrackDriver(context.Background(), driverID)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, driverID, res.DriverID)
		assert.Equal(t, -6.175392, res.Latitude)
		assert.Equal(t, 106.827153, res.Longitude)
		assert.Equal(t, ts, res.Timestamp)
	})

	t.Run("Error Case - driver_id tidak ditemukan", func(t *testing.T) {
		driverID := "driver-not-found"

		mockRepo.EXPECT().
			TrackDriver(gomock.Any(), gomock.Eq(driverID)).
			Return(nil, errors.New("driver location not found")).
			Times(1)

		res, err := svc.TrackDriver(context.Background(), driverID)
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, "driver location not found", err.Error())
	})
}
