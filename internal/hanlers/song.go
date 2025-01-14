package hanlers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"songs/internal/storages"
	"strconv"
)

// GetSongs
// @Summary Получить список песен
// @Description Получить список песен с пагинацией, отфильтрованный по названию группы и песни
// @Tags Песни
// @Param group query string false "Фильтр по названию группы"
// @Param song query string false "Фильтр по названию песни"
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество элементов на странице" default(10)
// @Success 200 {array} storages.Song
// @Failure 400 {object} map[string]interface{} "Неверные параметры запроса"
// @Failure 500 {object} map[string]interface{} "Не удалось получить песни"
// @Router /songs [get]
func (h *Handler) GetSongs(c *gin.Context) {
	group := c.Query("group")
	song := c.Query("song")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	h.logger.Infof("Получение песен с group=%s, song=%s, page=%d, limit=%d", group, song, page, limit)

	songs, err := h.storage.GetSongs(group, song, page, limit)
	if err != nil {
		h.logger.Errorf("Не удалось получить песни: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить песни", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, songs)
}

// GetLyrics
// @Summary Получить текст песни
// @Description Получить текст песни с пагинацией для конкретной песни
// @Tags Тексты песен
// @Param id path int true "ID песни"
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество элементов на странице" default(1)
// @Success 200 {array} string
// @Failure 400 {object} map[string]interface{} "Неверные параметры запроса"
// @Failure 500 {object} map[string]interface{} "Не удалось получить текст песни"
// @Router /songs/{id}/lyrics [get]
func (h *Handler) GetLyrics(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "1"))

	h.logger.Infof("Получение текста песни с ID=%d, page=%d, limit=%d", id, page, limit)

	lyrics, err := h.storage.GetLyrics(id, page, limit)
	if err != nil {
		h.logger.Errorf("Не удалось получить текст песни с ID=%d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить текст песни", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, lyrics)
}

// DeleteSong
// @Summary Удалить песню
// @Description Удалить песню по ее ID
// @Tags Песни
// @Param id path int true "ID песни"
// @Success 200 {object} map[string]interface{} "Песня удалена"
// @Failure 400 {object} map[string]interface{} "Неверный ID песни"
// @Failure 404 {object} map[string]interface{} "Песня не найдена"
// @Failure 500 {object} map[string]interface{} "Не удалось удалить песню"
// @Router /songs/{id} [delete]

func (h *Handler) DeleteSong(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	h.logger.Infof("Попытка удалить песню с ID=%d", id)

	err := h.storage.DeleteSong(id)
	if err != nil {
		h.logger.Errorf("Не удалось удалить песню с ID=%d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось удалить песню", "details": err.Error()})
		return
	}

	h.logger.Infof("Песня с ID=%d успешно удалена", id)
	c.JSON(http.StatusOK, gin.H{"message": "Песня удалена"})
}

// UpdateSongPartial
// @Summary Частичное обновление информации о песне
// @Description Обновить одно или несколько свойств песни по ее ID
// @Tags Песни
// @Param id path int true "ID песни"
// @Param song body storages.Song true "Данные о песне"
// @Success 200 {object} map[string]interface{} "Успешно"
// @Failure 400 {object} map[string]interface{} "Неверные данные для обновления песни"
// @Failure 404 {object} map[string]interface{} "Песня не найдена"
// @Failure 500 {object} map[string]interface{} "Не удалось обновить песню"
// @Router /songs/{id}/partial [put]
func (h *Handler) UpdateSongPartial(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		h.logger.Errorf("Неверные данные для обновления песни с ID=%d: %v", id, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные", "details": err.Error()})
		return
	}

	h.logger.Infof("Обновление песни с ID=%d, изменения=%+v", id, updates)

	err := h.storage.UpdateSongPartial(id, updates)
	if err != nil {
		h.logger.Errorf("Не удалось обновить песню с ID=%d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось обновить песню", "details": err.Error()})
		return
	}

	h.logger.Infof("Песня с ID=%d успешно обновлена", id)
	c.JSON(http.StatusOK, gin.H{"message": "Песня обновлена"})
}

// AddSong
// @Summary Добавить новую песню
// @Description Добавить новую песню в базу данных
// @Tags Песни
// @Param song body storages.Song true "Данные о песне"
// @Success 201 {object} map[string]interface{} "Песня добавлена"
// @Failure 400 {object} map[string]interface{} "Неверные данные для добавления песни"
// @Failure 500 {object} map[string]interface{} "Не удалось добавить песню"
// @Router /songs [post]
func (h *Handler) AddSong(c *gin.Context) {
	var song storages.Song
	if err := c.ShouldBindJSON(&song); err != nil {
		h.logger.Errorf("Неверные данные для добавления новой песни: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные", "details": err.Error()})
		return
	}

	h.logger.Infof("Добавление новой песни: %+v", song)

	songDetail, err := h.getSongDetail(song.Group, song.Name)
	if err != nil {
		h.logger.Errorf("Не удалось получить детали песни с внешнего API: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить детали песни", "details": err.Error()})
		return
	}

	h.logger.Infof("Детали песни с внешнего API: %+v", songDetail)

	song.ReleaseDate = songDetail.ReleaseDate
	song.Link = songDetail.Link

	err = h.storage.AddSong(song)
	if err != nil {
		h.logger.Errorf("Не удалось добавить песню: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось добавить песню", "details": err.Error()})
		return
	}

	for _, line := range songDetail.Text {
		err = h.storage.AddLyrics(song.ID, line)
		if err != nil {
			h.logger.Errorf("Не удалось добавить текст песни: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось добавить текст песни", "details": err.Error()})
			return
		}
	}

	h.logger.Infof("Песня успешно добавлена: %+v", song)
	c.JSON(http.StatusCreated, gin.H{"message": "Песня добавлена", "details": songDetail})
}

func (h *Handler) getSongDetail(group string, song string) (*storages.SongDetail, error) {
	apiURL := fmt.Sprintf("%s/info?group=%s&song=%s", h.config.ExternalAPI.Address, group, song)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("получен ответ с кодом отличным от 200 от внешнего API: %s", resp.Status)
	}

	var songDetail storages.SongDetail
	if err := json.NewDecoder(resp.Body).Decode(&songDetail); err != nil {
		return nil, fmt.Errorf("не удалось декодировать ответ: %v", err)
	}

	return &songDetail, nil
}
