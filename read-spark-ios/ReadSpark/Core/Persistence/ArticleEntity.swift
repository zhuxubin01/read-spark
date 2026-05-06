import CoreData
import Foundation

@objc(ArticleEntity)
public class ArticleEntity: NSManagedObject {
    @nonobjc public class func fetchRequest() -> NSFetchRequest<ArticleEntity> {
        NSFetchRequest<ArticleEntity>(entityName: "ArticleEntity")
    }

    @NSManaged public var id: UUID?
    @NSManaged public var title: String?
    @NSManaged public var summary: String?
    @NSManaged public var category: String?
    @NSManaged public var difficulty: String?
    @NSManaged public var wordCount: Int32
    @NSManaged public var isPremium: Bool
    @NSManaged public var updatedAt: Date?
}
