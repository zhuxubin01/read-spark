import Combine
import Foundation

protocol APIServiceProtocol {
    func login(phone: String, code: String) -> AnyPublisher<TokenPair, APIError>
    func register(phone: String, code: String) -> AnyPublisher<TokenPair, APIError>
    func syncProgress(articleId: UUID, position: Int, percentage: Double) -> AnyPublisher<[String: String], APIError>
}

final class APIService {
    static let shared = APIService()
    private let client = APIClient.shared

    private struct AuthRequest: Codable {
        let phone: String
        let code: String
    }

    private struct ProgressRequest: Codable {
        let article_id: String
        let position: Int
        let percentage: Double
    }

    // Auth
    func login(phone: String, code: String) -> AnyPublisher<TokenPair, APIError> {
        let body = try? JSONEncoder().encode(AuthRequest(phone: phone, code: code))
        return client.request("/auth/login", method: "POST", body: body, authorized: false)
    }

    func register(phone: String, code: String) -> AnyPublisher<TokenPair, APIError> {
        let body = try? JSONEncoder().encode(AuthRequest(phone: phone, code: code))
        return client.request("/auth/register", method: "POST", body: body, authorized: false)
    }

    // Articles
    func getDailyArticles() -> AnyPublisher<DailyArticlesResponse, APIError> {
        client.request("/articles/daily")
    }

    func getArticle(id: UUID) -> AnyPublisher<Article, APIError> {
        client.request("/articles/\(id.uuidString)")
    }

    func getArticles(category: String? = nil,
                     difficulty: String? = nil,
                     page: Int = 1) -> AnyPublisher<ArticleListResponse, APIError> {
        var components = URLComponents(string: client.baseURL + "/articles")
        var queryItems: [URLQueryItem] = [URLQueryItem(name: "page", value: String(page))]

        if let category {
            queryItems.append(URLQueryItem(name: "category", value: category))
        }
        if let difficulty {
            queryItems.append(URLQueryItem(name: "difficulty", value: difficulty))
        }

        components?.queryItems = queryItems

        guard let url = components?.url else {
            return Fail(error: APIError.invalidURL).eraseToAnyPublisher()
        }

        var request = URLRequest(url: url)
        request.httpMethod = "GET"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        if let token = KeychainManager.shared.accessToken {
            request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        }

        return URLSession.shared.dataTaskPublisher(for: request)
            .tryMap { data, response in
                guard let httpResponse = response as? HTTPURLResponse else {
                    throw APIError.invalidResponse
                }
                guard (200...299).contains(httpResponse.statusCode) else {
                    throw APIError.serverError(String(data: data, encoding: .utf8) ?? "Unknown error")
                }
                return data
            }
            .decode(type: ArticleListResponse.self, decoder: JSONDecoder())
            .mapError { error in
                if let apiError = error as? APIError {
                    return apiError
                }
                return APIError.decodingError
            }
            .receive(on: DispatchQueue.main)
            .eraseToAnyPublisher()
    }

    // Progress
    func syncProgress(articleId: UUID,
                      position: Int,
                      percentage: Double) -> AnyPublisher<[String: String], APIError> {
        let body = try? JSONEncoder().encode(
            ProgressRequest(article_id: articleId.uuidString, position: position, percentage: percentage)
        )

        return client.request("/progress", method: "POST", body: body)
    }
}

extension APIService: APIServiceProtocol {}

struct DailyArticlesResponse: Codable {
    let articles: [ArticleSummary]
}

struct ArticleListResponse: Codable {
    let articles: [ArticleSummary]
    let total: Int
    let page: Int
    let page_size: Int
}
