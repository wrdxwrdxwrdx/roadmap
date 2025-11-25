import axios, { AxiosInstance, InternalAxiosRequestConfig, AxiosResponse, AxiosError } from 'axios'

// Создание экземпляра axios с базовой конфигурацией
const apiClient: AxiosInstance = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Интерсептор для запросов - добавляет токены авторизации, логирование и т.д.
apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    // Здесь можно добавить токен авторизации из store или localStorage
    // const token = localStorage.getItem('token')
    // if (token && config.headers) {
    //   config.headers.Authorization = `Bearer ${token}`
    // }
    
    // Логирование запросов (только в development)
    if (import.meta.env.DEV) {
      console.log(`[API Request] ${config.method?.toUpperCase()} ${config.url}`, config)
    }
    
    return config
  },
  (error: AxiosError) => {
    console.error('[API Request Error]', error)
    return Promise.reject(error)
  }
)

// Интерсептор для ответов - обработка ошибок, логирование
apiClient.interceptors.response.use(
  (response: AxiosResponse) => {
    // Логирование ответов (только в development)
    if (import.meta.env.DEV) {
      console.log(`[API Response] ${response.config.method?.toUpperCase()} ${response.config.url}`, response.data)
    }
    
    return response
  },
  (error: AxiosError) => {
    // Обработка ошибок
    if (error.response) {
      // Сервер ответил с кодом ошибки
      switch (error.response.status) {
        case 401:
          // Не авторизован - можно перенаправить на страницу логина
          console.error('[API Error] Unauthorized')
          // window.location.href = '/login'
          break
        case 403:
          console.error('[API Error] Forbidden')
          break
        case 404:
          console.error('[API Error] Not Found')
          break
        case 500:
          console.error('[API Error] Server Error')
          break
        default:
          console.error('[API Error]', error.response.status, error.response.data)
      }
    } else if (error.request) {
      // Запрос был сделан, но ответа не получено
      console.error('[API Error] No response received', error.request)
    } else {
      // Ошибка при настройке запроса
      console.error('[API Error]', error.message)
    }
    
    return Promise.reject(error)
  }
)

export default apiClient

