import Foundation

struct User: Codable, Identifiable {
    let id: UUID
    let phone: String
    let email: String?
    let nickname: String?
    let avatarUrl: String?
    let createdAt: String
}
