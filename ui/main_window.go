package ui

import (
	"fmt"
	"strings"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"gorm.io/gorm"
	"szi-registry/export"
	"szi-registry/models"
	"szi-registry/services"
	"szi-registry/utils/ui"
	"time"
)

// ShowMainWindow отображает главное окно приложения
func ShowMainWindow(myApp fyne.App, username string, db *gorm.DB) {
	myWindow := myApp.NewWindow("Реестр СЗИ от НСД — " + username)
	myWindow.Resize(fyne.NewSize(1400, 800))

	// Получаем ID пользователя по имени
	user, err := services.GetUserByUsername(db, username)
	if err != nil || user == nil {
		// Если не удалось получить пользователя, используем 0
		user = &models.User{ID: 0}
	}

	// Создаем срезы для хранения данных и действий
	var tableData [][]string
	var buttonActions [][2]func()
	var table *widget.Table

	// Загрузка данных из базы данных для конкретного пользователя
	records, err := services.GetUserRecords(db, uint(user.ID))
	if err != nil {
		// В реальном приложении нужно обработать ошибку
		records = []models.SZIRecord{} // Используем пустой массив
	}

	// Функция для обновления данных в таблице
	updateTableData := func(data []models.SZIRecord) {
		tableData = make([][]string, 0)
		buttonActions = make([][2]func(), 0)

		for _, record := range data {
			// Пропускаем записи с пустым названием
			if record.Name == "" {
				continue
			}

			// Определяем статус на основе даты истечения
			status := record.Status

			// Если статус не установлен явно, определяем его по дате истечения
			if status == "" {
				status = ui.GetRecordStatus(record)
			}

			rowData := []string{
				record.Name,
				record.Type,
				record.CertNumber,
				record.CertIssueDate.Format("2006-01-02"),
				record.CertExpiryDate.Format("2006-01-02"),
				status,
				record.Manufacturer,
				record.SoftwareVersion,
				record.Purpose,
				record.DeploymentType,
				record.ClassProtection,
			}

			tableData = append(tableData, rowData)

			// Сохраняем действия для каждой строки
			currentRecord := record // захватываем значение для замыкания
			var confirmDialog *widget.PopUp
			actions := [2]func(){
				func() { // Редактировать
					ShowEditSziWindow(myApp, currentRecord, db)
				},
				func() { // Удалить
					// Создаем диалог подтверждения вручную
					confirmDialog = widget.NewModalPopUp(
						container.NewVBox(
							widget.NewLabel("Вы уверены, что хотите удалить запись '"+currentRecord.Name+"'?"),
							container.NewHBox(
								widget.NewButton("Да", func() {
									err := services.DeleteSziRecord(db, uint(user.ID), currentRecord.ID)
									if err != nil {
										// Показать ошибку
										ui.ShowErrorDialog("Ошибка при удалении записи: "+err.Error(), myWindow.Canvas())
									} else {
										// Обновить окно
										myWindow.Close()
										ShowMainWindow(myApp, username, db)
									}
									confirmDialog.Hide()
								}),
								widget.NewButton("Нет", func() {
									confirmDialog.Hide()
								}),
							),
						),
						myWindow.Canvas(),
					)
					confirmDialog.Show()
				},
			}
			buttonActions = append(buttonActions, actions)
		}
	}

	// Инициализируем данные
	updateTableData(records)

	// Определяем начальные ширины столбцов
	initialColumnWidths := []float32{250, 120, 140, 120, 120, 100, 150, 120, 120, 120, 120, 100}

	// Создаем таблицу с 12 столбцами (11 данных + 1 для действий)
	table = widget.NewTable(
		func() (int, int) {
			return len(tableData), 12 // строки: записи из БД, 12 столбцов
		},
		func() fyne.CanvasObject {
			// Возвращаем метку для обычных ячеек
			label := widget.NewLabel("")
			return label
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)

			if id.Row < len(tableData) {
				if id.Col < 11 {
					// Обычные данные
					if id.Col < len(tableData[id.Row]) {
						text := tableData[id.Row][id.Col]
						label.SetText(text)
					} else {
						label.SetText("")
					}
				} else if id.Col == 11 {
					// Столбец действий
					label.SetText("...")
					label.Alignment = fyne.TextAlignCenter
				}
			} else {
				label.SetText("")
			}
		})

	// Устанавливаем начальные ширины столбцов
	for i, width := range initialColumnWidths {
		if i < len(initialColumnWidths) {
			table.SetColumnWidth(i, width)
		}
	}

	// Обработка нажатий на таблицу
	table.OnSelected = func(id widget.TableCellID) {
		if id.Col == 11 && id.Row < len(buttonActions) {
			// Показываем диалог выбора действия
			actionMenu := widget.NewPopUpMenu(
				&fyne.Menu{
					Label: "Действия",
					Items: []*fyne.MenuItem{
						{Label: "Редактировать", Action: buttonActions[id.Row][0]},
						{Label: "Удалить", Action: buttonActions[id.Row][1]},
					},
				}, myWindow.Canvas())

			// Показываем меню в позиции курсора
			actionMenu.Show()
		}
	}

	// Заголовки столбцов
	headers := []string{"Наименование СЗИ", "Тип СЗИ", "№ сертификата", "Дата выдачи", "Срок действия", "Статус", "Производитель", "Версия ПО", "Назначение", "Тип развертывания", "Класс защищенности", "Действия"}
	headerContainer := container.NewHBox()
	for _, header := range headers {
		label := widget.NewLabelWithStyle(header, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
		label.Alignment = fyne.TextAlignCenter
		box := container.NewVBox(label, widget.NewSeparator())
		headerContainer.Add(box)
	}

	// Создаем поля для фильтрации
	nameFilter := widget.NewEntry()
	nameFilter.SetPlaceHolder("Поиск по наименованию...")
	typeFilter := widget.NewSelectEntry([]string{"СКЗИ", "ОС", "СЗИ КС", "СЗИ СКЗИ", "Другое"})
	typeFilter.SetPlaceHolder("Тип СЗИ...")
	statusFilter := widget.NewSelectEntry([]string{"Актуально", "Просрочено", "Снято"})
	statusFilter.SetPlaceHolder("Статус...")
	locationFilter := widget.NewEntry()
	locationFilter.SetPlaceHolder("Место установки...")
	manufacturerFilter := widget.NewEntry()
	manufacturerFilter.SetPlaceHolder("Производитель...")
	purposeFilter := widget.NewSelectEntry([]string{"АС", "ИВК", "Универсальное"})
	purposeFilter.SetPlaceHolder("Назначение...")
	deploymentTypeFilter := widget.NewSelectEntry([]string{"Клиент-сервер", "Автономное", "АПК", "Виртуальное"})
	deploymentTypeFilter.SetPlaceHolder("Тип развертывания...")
	classProtectionFilter := widget.NewSelectEntry([]string{"1", "2", "3А", "4", "5"})
	classProtectionFilter.SetPlaceHolder("Класс защищенности...")

	// Функция для применения фильтров
	applyFilters := func() {
		filters := make(map[string]interface{})

		if nameFilter.Text != "" {
			filters["name"] = nameFilter.Text
		}
		if typeFilter.Text != "" {
			filters["type"] = typeFilter.Text
		}
		if statusFilter.Text != "" {
			filters["status"] = statusFilter.Text
		}
		if locationFilter.Text != "" {
			filters["location"] = locationFilter.Text
		}
		if manufacturerFilter.Text != "" {
			filters["manufacturer"] = manufacturerFilter.Text
		}
		if purposeFilter.Text != "" {
			filters["purpose"] = purposeFilter.Text
		}
		if deploymentTypeFilter.Text != "" {
			filters["deployment_type"] = deploymentTypeFilter.Text
		}
		if classProtectionFilter.Text != "" {
			filters["class_protection"] = classProtectionFilter.Text
		}

		filteredRecords, err := services.SearchSZIRecords(db, uint(user.ID), filters)
		if err != nil {
			// Показать ошибку
			ui.ShowErrorDialog("Ошибка при поиске записей: "+err.Error(), myWindow.Canvas())
			return
		}

		// Обновляем данные в таблице
		updateTableData(filteredRecords)
		table.Refresh()
	}

	// Кнопка сброса фильтров
	resetFilters := func() {
		nameFilter.SetText("")
		typeFilter.SetText("")
		statusFilter.SetText("")
		locationFilter.SetText("")
		manufacturerFilter.SetText("")
		purposeFilter.SetText("")
		deploymentTypeFilter.SetText("")
		classProtectionFilter.SetText("")

		// Загружаем все записи снова
		allRecords, err := services.GetUserRecords(db, uint(user.ID))
		if err != nil {
			// Показать ошибку
			ui.ShowErrorDialog("Ошибка при загрузке записей: "+err.Error(), myWindow.Canvas())
			return
		}

		// Обновляем данные в таблице
		updateTableData(allRecords)
		table.Refresh()
	}

	// Контейнер для фильтров
	filterContainer := container.NewGridWithColumns(7,
		container.NewVBox(widget.NewLabel("Поиск по наименованию:"), nameFilter),
		container.NewVBox(widget.NewLabel("Тип СЗИ:"), typeFilter),
		container.NewVBox(widget.NewLabel("Статус:"), statusFilter),
		container.NewVBox(widget.NewLabel("Место установки:"), locationFilter),
		container.NewVBox(widget.NewLabel("Производитель:"), manufacturerFilter),
		container.NewVBox(widget.NewLabel("Назначение:"), purposeFilter),
		container.NewVBox(
			widget.NewButton("Применить фильтры", applyFilters),
			widget.NewButton("Сбросить фильтры", resetFilters),
		),
	)

	// Добавляем кнопки для межпользовательского взаимодействия
	shareBtn := widget.NewButton("Поделиться записью", func() {
		// Открываем диалог для выбора записи и пользователя для предоставления доступа
		dialog := createShareDialog(myWindow, uint(user.ID), db)
		dialog.Show()
	})

	requestAccessBtn := widget.NewButton("Запросить доступ", func() {
		// Открываем диалог для запроса доступа к чужой записи
		dialog := createRequestAccessDialog(myWindow, uint(user.ID), db)
		dialog.Show()
	})

	notificationsBtn := widget.NewButton("Уведомления", func() {
		// Открываем окно с уведомлениями
		dialog := createNotificationsDialog(myWindow, uint(user.ID), db)
		dialog.Show()
	})


	// Кнопка импорта из CSV
	importBtn := widget.NewButton("Импорт из CSV", func() {
		// Диалог для выбора CSV файла
		fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog := widget.NewModalPopUp(
					widget.NewLabel("Ошибка при открытии файла: "+err.Error()),
					myWindow.Canvas(),
				)
				dialog.Show()
				return
			}

			if reader == nil {
				return // Пользователь отменил выбор файла
			}

			// Получаем путь к файлу
			filePath := reader.URI().String()

			// Вызываем сервис импорта
			importedCount, err := services.ImportSziRecordsFromCSV(db, uint(user.ID), filePath)
			if err != nil {
				dialog := widget.NewModalPopUp(
					widget.NewLabel("Ошибка при импорте из CSV: "+err.Error()),
					myWindow.Canvas(),
				)
				dialog.Show()
			} else {
				// Обновляем данные в таблице
				records, err := services.GetUserRecords(db, uint(user.ID))
				if err != nil {
					// В реальном приложении нужно обработать ошибку
					records = []models.SZIRecord{} // Используем пустой массив
				}

				// Обновляем данные таблицы
				updateTableData(records)
				table.Refresh()

				// Показываем сообщение об успешном импорте
				dialog := widget.NewModalPopUp(
					widget.NewLabel(fmt.Sprintf("Успешно импортировано %d записей из CSV", importedCount)),
					myWindow.Canvas(),
				)
				dialog.Show()
			}
		}, myWindow)

		fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".csv"}))
		fileDialog.Show()
	})

	// Кнопка экспорта в CSV
	exportBtn := widget.NewButton("Экспорт в CSV", func() {
		// Получаем все записи пользователя для экспорта
		records, err := services.GetUserRecords(db, uint(user.ID))
		if err != nil {
			dialog := widget.NewModalPopUp(
				widget.NewLabel("Ошибка при загрузке данных для экспорта: "+err.Error()),
				myWindow.Canvas(),
			)
			dialog.Show()
			return
		}

		// Генерируем имя файла с временной меткой
		filename := fmt.Sprintf("szi_export_%s.csv", time.Now().Format("20060102_150405"))

		// Экспортируем данные
		err = export.ExportToCSV(records, filename)
		if err != nil {
			dialog := widget.NewModalPopUp(
				widget.NewLabel("Ошибка при экспорте в CSV: "+err.Error()),
				myWindow.Canvas(),
			)
			dialog.Show()
		} else {
			dialog := widget.NewModalPopUp(
				widget.NewLabel("Данные успешно экспортированы в файл: "+filename),
				myWindow.Canvas(),
			)
			dialog.Show()
		}
	})

	addBtn := widget.NewButton("Добавить СЗИ", func() {
		user, err := services.GetUserByUsername(db, username)
		if err != nil || user == nil {
			// Если не удалось получить пользователя, показываем ошибку
			dialog := widget.NewModalPopUp(
				widget.NewLabel("Не удалось получить информацию о пользователе"),
				myWindow.Canvas(),
			)
			dialog.Show()
			return
		}
		ShowAddSziWindow(myApp, uint(user.ID), db)
	})

	logoutBtn := widget.NewButton("Выход", func() {
		myWindow.Close()
		ShowLoginWindow(myApp, db)
	})

	// Создаем контейнер с заголовками и таблицей, чтобы они прокручивались вместе
	tableWithHeaders := container.NewVBox(
		headerContainer,
		table,
	)

	// Создаем контейнер с правильным расположением элементов
	content := container.NewBorder(
		filterContainer,                                     // верхняя часть - фильтры
		container.NewHBox(logoutBtn, addBtn),                // нижняя часть - основные кнопки
		nil,                                                 // левая часть - нет
		nil,                                                 // правая часть - нет
		tableWithHeaders,                                    // центральная часть - заголовки и таблица
	)

	// Добавляем панель инструментов для межпользовательского взаимодействия
	toolbar := container.NewHBox(
		shareBtn,
		requestAccessBtn,
		notificationsBtn,
		importBtn,
		exportBtn,
		widget.NewButton("Статистика", func() {
			showStatisticsWindow(myApp, uint(user.ID), db)
		}),
	)

	// Объединяем основной контент с панелью инструментов
	finalContent := container.NewBorder(toolbar, nil, nil, nil, content)

	// Обернем в прокрутку для больших таблиц
	scrollableContent := container.NewScroll(finalContent)
	scrollableContent.SetMinSize(fyne.NewSize(1300, 700))

	myWindow.SetContent(scrollableContent)
	myWindow.Show()
}

// createShareDialog создает диалог для предоставления доступа к записи другому пользователю
func createShareDialog(parentWindow fyne.Window, currentUserID uint, db *gorm.DB) *widget.PopUp {
	// Создаем элементы интерфейса для выбора записи и пользователя
	recordSelector := widget.NewEntry()
	recordSelector.SetPlaceHolder("Введите название СЗИ")
	userSelector := widget.NewSelectEntry([]string{})
	permissionSelector := widget.NewSelect([]string{"r", "rw"}, func(string) {})

	// Создаем контейнер для предложений записей
	suggestionsContainer := container.NewVBox()

	// Создаем диалог
	var dialog *widget.PopUp

	// Создаем списки записей и пользователей
	var allRecords []models.SZIRecord
	var filteredRecords []models.SZIRecord

	// Загружаем доступные записи и пользователей
	loadData := func() {
		// Загружаем записи пользователя
		records, err := services.GetUserRecords(db, currentUserID)
		if err != nil {
			// Обработка ошибки
			fmt.Printf("Ошибка при загрузке записей: %v\n", err)
			return
		}

		allRecords = records
		filteredRecords = records

		// Загружаем всех пользователей, кроме текущего
		allUsers, err := services.GetAllUsers(db)
		if err != nil {
			// Обработка ошибки
			fmt.Printf("Ошибка при загрузке пользователей: %v\n", err)
			return
		}

		var userOptions []string
		for _, user := range allUsers {
			if user.ID != currentUserID { // Исключаем текущего пользователя
				userOptions = append(userOptions, user.Username)
			}
		}
		userSelector.SetOptions(userOptions)
	}

	// Функция для обновления списка записей на основе введенного текста
	updateSuggestions := func() {
		// Очищаем предыдущие предложения
		suggestionsContainer.Objects = nil

		// Фильтруем записи по введенному тексту
		inputText := recordSelector.Text
		filteredRecords = []models.SZIRecord{}

		for _, record := range allRecords {
			if containsSubstring(record.Name, inputText) {
				filteredRecords = append(filteredRecords, record)
			}
		}

		// Создаем кнопки для каждой подходящей записи
		for _, record := range filteredRecords {
			btn := widget.NewButton(record.Name, func() {
				// При выборе записи, устанавливаем её название в поле ввода
				recordSelector.SetText(record.Name)
				// Скрываем список предложений
				suggestionsContainer.Hide()
			})
			suggestionsContainer.Add(btn)
		}

		// Обновляем интерфейс
		if len(filteredRecords) > 0 {
			suggestionsContainer.Show()
		} else {
			suggestionsContainer.Hide()
		}
		suggestionsContainer.Refresh()
	}

	loadData()

	// Обработчик изменения текста в поле ввода
	recordSelector.OnChanged = func(text string) {
		updateSuggestions()
	}

	// Кнопка подтверждения
	confirmBtn := widget.NewButton("Поделиться", func() {
		// Логика предоставления доступа
		// Нужно найти ID выбранной записи
		var selectedRecordID uint
		var found bool

		// Сначала ищем в исходном списке
		for _, record := range allRecords {
			if record.Name == recordSelector.Text {
				selectedRecordID = record.ID
				found = true
				break
			}
		}

		// Если не нашли, ищем в отфильтрованном списке
		if !found {
			for _, record := range filteredRecords {
				if record.Name == recordSelector.Text {
					selectedRecordID = record.ID
					found = true
					break
				}
			}
		}

		if !found {
			// Показать ошибку
			errorContent := container.NewVBox(
				widget.NewLabel("Выбранная запись не найдена"),
				widget.NewButton("OK", func() {
					dialog.Hide()
				}),
			)
			errorDialog := widget.NewModalPopUp(errorContent, parentWindow.Canvas())
			errorDialog.Show()
			return
		}

		// Получаем ID выбранного пользователя
		allUsers, err := services.GetAllUsers(db)
		if err != nil {
			// Показать ошибку
			ui.ShowErrorDialog("Ошибка при загрузке пользователей: "+err.Error(), parentWindow.Canvas())
			return
		}

		var selectedUserID uint
		found = false
		for _, user := range allUsers {
			if user.Username == userSelector.Text && user.ID != currentUserID {
				selectedUserID = user.ID
				found = true
				break
			}
		}

		if !found {
			ui.ShowErrorDialog("Выбранный пользователь не найден", parentWindow.Canvas())
			return
		}

		// Вызываем сервис для предоставления доступа
		err = services.GrantAccess(db, currentUserID, selectedUserID, selectedRecordID, permissionSelector.Selected, nil)
		if err != nil {
			// Показать ошибку
			ui.ShowErrorDialog("Ошибка при предоставлении доступа: "+err.Error(), parentWindow.Canvas())
		} else {
			// Показать успех
			ui.ShowSuccessDialog("Доступ успешно предоставлен!", parentWindow.Canvas())
		}
	})

	// Кнопка отмены
	cancelBtn := widget.NewButton("Отмена", func() {
		dialog.Hide()
	})

	dialogContent := container.NewVBox(
		widget.NewLabel("Название СЗИ:"),
		recordSelector,
		suggestionsContainer,
		widget.NewLabel("Выберите пользователя:"),
		userSelector,
		widget.NewLabel("Выберите права доступа:"),
		permissionSelector,
		container.NewHBox(confirmBtn, cancelBtn),
	)

	dialog = widget.NewModalPopUp(dialogContent, parentWindow.Canvas())
	return dialog
}

// создает диалог для запроса доступа к чужой записи
func createRequestAccessDialog(parentWindow fyne.Window, currentUserID uint, db *gorm.DB) *widget.PopUp {
	var dialog *widget.PopUp

	// Создаем элементы интерфейса для запроса доступа
	recordSelector := widget.NewEntry()
	recordSelector.SetPlaceHolder("Введите название СЗИ")
	permissionSelector := widget.NewSelect([]string{"r", "rw"}, func(string) {})
	messageEntry := widget.NewMultiLineEntry()
	messageEntry.SetPlaceHolder("Сообщение владельцу")

	// Создаем контейнер для предложений
	suggestionsContainer := container.NewVBox()

	// Создаем список записей для отображения
	var allRecords []models.SZIRecord
	var filteredRecords []models.SZIRecord

	// Загружаем все доступные записи
	loadRecords := func() {
		// Получаем все записи, к которым у пользователя нет доступа
		records, err := services.GetSziRecordsAccessibleByUserID(db, currentUserID)
		if err != nil {
			// Обработка ошибки
			fmt.Printf("Ошибка при загрузке записей: %v\n", err)
			return
		}

		// Получаем уже имеющиеся разрешения
		permissions, err := services.GetUserPermissions(db, currentUserID)
		if err != nil {
			// Обработка ошибки
			fmt.Printf("Ошибка при загрузке разрешений: %v\n", err)
			return
		}

		// Создаем множество ID записей, к которым уже есть доступ
		accessibleRecordIDs := make(map[uint]bool)
		for _, perm := range permissions {
			accessibleRecordIDs[perm.RecordID] = true
		}

		// Фильтруем записи, к которым нет доступа
		allRecords = []models.SZIRecord{}
		for _, record := range records {
			// Добавляем только те записи, к которым нет доступа
			if !accessibleRecordIDs[record.ID] {
				allRecords = append(allRecords, record)
			}
		}

		// Изначально все записи входят в фильтрованный список
		filteredRecords = allRecords
	}

	// Функция для обновления списка записей на основе введенного текста
	updateSuggestions := func() {
		// Очищаем предыдущие предложения
		suggestionsContainer.Objects = nil

		// Фильтруем записи по введенному тексту
		inputText := recordSelector.Text
		filteredRecords = []models.SZIRecord{}

		for _, record := range allRecords {
			if containsSubstring(record.Name, inputText) {
				filteredRecords = append(filteredRecords, record)
			}
		}

		// Создаем кнопки для каждой подходящей записи
		for _, record := range filteredRecords {
			btn := widget.NewButton(record.Name, func() {
				// При выборе записи, устанавливаем её название в поле ввода
				recordSelector.SetText(record.Name)
				// Скрываем список предложений
				suggestionsContainer.Hide()
			})
			suggestionsContainer.Add(btn)
		}

		// Обновляем интерфейс
		if len(filteredRecords) > 0 {
			suggestionsContainer.Show()
		} else {
			suggestionsContainer.Hide()
		}
		suggestionsContainer.Refresh()
	}

	// Загружаем записи при открытии диалога
	loadRecords()

	// Обработчик изменения текста в поле ввода
	recordSelector.OnChanged = func(text string) {
		updateSuggestions()
	}

	// Кнопка подтверждения
	confirmBtn := widget.NewButton("Запросить доступ", func() {
		// Находим ID записи по названию
		var selectedRecord *models.SZIRecord
		for _, record := range allRecords {
			if record.Name == recordSelector.Text {
				selectedRecord = &record
				break
			}
		}

		if selectedRecord == nil {
			// Проверяем также среди отфильтрованных записей
			for _, record := range filteredRecords {
				if record.Name == recordSelector.Text {
					selectedRecord = &record
					break
				}
			}
		}

		if selectedRecord == nil {
			// Показать ошибку - запись не найдена
			ui.ShowErrorDialog("Выбранная запись не найдена", parentWindow.Canvas())
			return
		}

		// Получаем владельца записи
		ownerID := selectedRecord.UserID

		// Вызываем сервис для запроса доступа
		err := services.RequestAccess(db, currentUserID, ownerID, selectedRecord.ID, permissionSelector.Selected, messageEntry.Text)
		if err != nil {
			// Показать ошибку
			errorDialog := widget.NewModalPopUp(
				widget.NewLabel("Ошибка при отправке запроса доступа: "+err.Error()),
				parentWindow.Canvas(),
			)
			errorDialog.Show()
		} else {
			// Показать успех
			successDialog := widget.NewModalPopUp(
				widget.NewLabel("Запрос доступа успешно отправлен!"),
				parentWindow.Canvas(),
			)
			successDialog.Show()

			// Закрываем диалог
			dialog.Hide()
		}
	})

	cancelBtn := widget.NewButton("Отмена", func() {
		dialog.Hide()
	})

	dialogContent := container.NewVBox(
		widget.NewLabel("Название СЗИ:"),
		recordSelector,
		suggestionsContainer,
		widget.NewLabel("Требуемые права:"),
		permissionSelector,
		widget.NewLabel("Сообщение владельцу:"),
		messageEntry,
		container.NewHBox(confirmBtn, cancelBtn),
	)

	dialog = widget.NewModalPopUp(dialogContent, parentWindow.Canvas())
	return dialog
}

// вспомогательная функция для проверки наличия подстроки в строке (без учета регистра)
func containsSubstring(str, substr string) bool {
	return strings.Contains(strings.ToLower(str), strings.ToLower(substr))
}

// создает диалог для просмотра уведомлений
func createNotificationsDialog(parentWindow fyne.Window, userID uint, db *gorm.DB) *widget.PopUp {
	var dialog *widget.PopUp

	// Загружаем уведомления пользователя
	notifications, err := services.GetNotifications(db, userID)
	if err != nil {
		// Обработка ошибки
		notifications = []models.Notification{}
	}

	var content fyne.CanvasObject

	if len(notifications) == 0 {
		// Показываем сообщение, если уведомлений нет
		noNotificationsLabel := widget.NewLabel("Уведомлений нет")
		noNotificationsLabel.Alignment = fyne.TextAlignCenter
		content = container.NewPadded(noNotificationsLabel)
	} else {
		// Создаем список уведомлений
		list := widget.NewList(
			func() int {
				return len(notifications)
			},
			func() fyne.CanvasObject {
				return container.NewVBox(
					widget.NewLabel("Title"),
					widget.NewLabel("Message"),
					widget.NewLabel("Details"),
					widget.NewButton("Пометить как прочитанное", func() {}),
				)
			},
			func(id widget.ListItemID, obj fyne.CanvasObject) {
				cont := obj.(*fyne.Container)
				labels := cont.Objects[:3]
				button := cont.Objects[3].(*widget.Button)

				titleLabel := labels[0].(*widget.Label)
				messageLabel := labels[1].(*widget.Label)
				detailsLabel := labels[2].(*widget.Label)

				notification := notifications[id]

				status := "непрочитано"
				if notification.ReadStatus {
					status = "прочитано"
				}

				titleLabel.SetText("[" + notification.Type + "] " + notification.Title)
				messageLabel.SetText(notification.Message)
				detailsLabel.SetText("От: " + fmt.Sprintf("%d", notification.SenderID) + ", Статус: " + status)

				button.SetText("Пометить как прочитанное")
				button.Hidden = notification.ReadStatus // Скрываем кнопку, если уже прочитано

				button.OnTapped = func() {
					// Помечаем уведомление как прочитанное
					err := services.MarkNotificationAsRead(db, uint(notification.ID))
					if err != nil {
						// Показать ошибку
						errorDialog := widget.NewModalPopUp(
							widget.NewLabel("Ошибка при пометке уведомления: "+err.Error()),
							parentWindow.Canvas(),
						)
						errorDialog.Show()
					} else {
						// Обновляем список
						parentWindow.Canvas().Content().Refresh()
					}
				}
			},
		)
		content = list
	}

	closeBtn := widget.NewButton("Закрыть", func() {
		dialog.Hide()
	})

	dialogContent := container.NewBorder(nil, closeBtn, nil, nil, content)
	dialog = widget.NewModalPopUp(dialogContent, parentWindow.Canvas())

	// Устанавливаем минимальный размер окна для большей ширины
	dialog.Resize(fyne.NewSize(600, 400))

	return dialog
}
