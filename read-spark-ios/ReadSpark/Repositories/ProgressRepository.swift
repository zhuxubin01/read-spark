import Combine
import Foundation

final class ProgressRepository {
    static let shared = ProgressRepository()
    private let api = APIService.shared

    func syncProgress(articleId: UUID, position: Int, percentage: Double) -> AnyPublisher<Void, APIError> {
        api.syncProgress(articleId: articleId, position: position, percentage: percentage)
            .map { _ in () }
            .eraseToAnyPublisher()
    }
}
